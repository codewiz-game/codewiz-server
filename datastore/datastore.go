package datastore

import (
	"database/sql"
	"errors"
	"github.com/go-gorp/gorp"
	_ "github.com/mattes/migrate/driver/mysql"
	_ "github.com/mattes/migrate/driver/sqlite3"
	"github.com/mattes/migrate/migrate"
	"reflect"
	"time"
	"path/filepath"
)

type StatusCode uint

const (
	Transient StatusCode = iota // not yet inserted - the zero-value for StatusCode
	Active
	Deleted

)

type DB struct {
	*gorp.DbMap
	driver string
	dsn    string
}

func Open(driver string, dsn string) (*DB, error) {
	db, err := sql.Open(driver, dsn)
	if err != nil {
		return nil, err
	}

	// Check that the connection is valid before continuing
	err = db.Ping()
	if err != nil {
		return nil, err
	}

	gorpDialect := getDialectForDriver(driver)
	dbMap := &gorp.DbMap{Db: db, Dialect: gorpDialect}
	return &DB{DbMap: dbMap, driver: driver, dsn: dsn}, nil
}

func (ds *DB) UpSync(migrationsPath string) ([]error, bool) {
	return migrate.UpSync(ds.driver+"://"+ds.dsn, filepath.Join(migrationsPath, ds.driver))
}

func (ds *DB) Insert(record interface{}) error {

	// Update the status, creation time and last modifed before insertion, so that the changes can be persisted
	statusRecord, hasStatus := record.(StatusRecorder)
	var previousStatus StatusCode
	if hasStatus {
		previousStatus = statusRecord.Status()
		statusRecord.SetStatus(Active)
	}

	now := getCurrentTime()
	creationTimeRecord, hasCreationTime := record.(CreationTimeRecorder)
	var previousCreationTime time.Time
	if hasCreationTime {
		previousCreationTime = creationTimeRecord.CreationTime()
		creationTimeRecord.SetCreationTime(now)
	}

	lastUpdatedTimeRecord, hasLastUpdatedTime := record.(LastUpdateTimeRecorder)
	var previousLastUpdatedTime time.Time
	if hasLastUpdatedTime {
		previousLastUpdatedTime = lastUpdatedTimeRecord.LastUpdatedTime()
		lastUpdatedTimeRecord.SetLastUpdatedTime(now)
	}

	var err error
	logicallyDeletableRecord, supportsLogicalDeletion := record.(LogicallyDeletable)
	if supportsLogicalDeletion {
		var existingRecord interface{}
		existingRecord, err = ds.DbMap.Get(record, logicallyDeletableRecord.Keys()...)
		if existingRecord != nil {
			// No need to check validity, since all LogicallyDeletable records are also StatusRecorders
			existingStatusRecorder := existingRecord.(StatusRecorder)
			if existingStatusRecorder.Status() == Deleted {
				err = ds.Update(record)
			} else {
				err = errors.New("A record with the same primary key already exists in the data store.")
			}
		} else {
			err = ds.DbMap.Insert(record)
		}
	} else {
		err = ds.DbMap.Insert(record)
	}

	if err != nil {
		// If the insertion failed, revert the changes made to
		// the status, creation and last modified time to keep the data accurate
		if hasStatus {
			statusRecord.SetStatus(previousStatus)
		}

		if hasCreationTime {
			creationTimeRecord.SetCreationTime(previousCreationTime)
		}

		if hasLastUpdatedTime {
			lastUpdatedTimeRecord.SetLastUpdatedTime(previousLastUpdatedTime)
		}
	}

	return err
}

func (ds *DB) Update(record interface{}) error {
	now := getCurrentTime()
	lastUpdatedTimeRecord, hasLastUpdatedTime := record.(LastUpdateTimeRecorder)
	var previousLastUpdatedTime time.Time
	if hasLastUpdatedTime {
		previousLastUpdatedTime = lastUpdatedTimeRecord.LastUpdatedTime()
		lastUpdatedTimeRecord.SetLastUpdatedTime(now)
	}

	count, err := ds.DbMap.Update(record)
	if err == nil && count == 0 {
		err = errors.New("No records were affected by the update operation.")
	}

	if err != nil {
		if hasLastUpdatedTime {
			lastUpdatedTimeRecord.SetLastUpdatedTime(previousLastUpdatedTime)
		}
	}

	return err
}

func (ds *DB) Delete(record interface{}) error {
	statusRecord, hasStatus := record.(StatusRecorder)
	var previousStatus StatusCode
	if hasStatus {
		previousStatus = statusRecord.Status()
		statusRecord.SetStatus(Deleted)
	}

	now := getCurrentTime()
	deletionTimeRecord, hasDeletionTime := record.(DeletionTimeRecorder)
	var previousDeletionTime time.Time
	if hasDeletionTime {
		previousDeletionTime = deletionTimeRecord.DeletionTime()
		deletionTimeRecord.SetDeletionTime(now)
	}

	err := ds.Update(record)
	if err != nil {
		if hasStatus {
			statusRecord.SetStatus(previousStatus)
		}

		if hasDeletionTime {
			deletionTimeRecord.SetDeletionTime(previousDeletionTime)
		}
	}

	return err
}

func (ds *DB) Get(record interface{}, whereClause string, args ...interface{}) (interface{}, error) {
	results, err := ds.Select(record, whereClause, args...)
	if err != nil {
		return nil, err
	}

	numResults := len(results)
	if numResults == 0 {
		return nil, nil
	}

	if numResults != 1 {
		return nil, errors.New("More than one matching record was found for the given query")
	}

	return results[0], nil
}

func (ds *DB) Select(results interface{}, whereClause string, args ...interface{}) ([]interface{}, error) {

	recordType := reflect.TypeOf(results)

	// Determine the record type by digging down until we reach something that isn't a container or reference
	for recordType.Kind() == reflect.Ptr || recordType.Kind() == reflect.Slice {
		recordType = recordType.Elem()
	}

	// Create a new pointer to the determined type
	recordPtr := reflect.New(recordType).Interface()

	// Use the pointer to determine whether the type supports logical deletion, and modify the query if it does
	logicallyDeletableRecord, supportsLogicalDeletion := recordPtr.(LogicallyDeletable)
	if supportsLogicalDeletion {
		statusColumn := logicallyDeletableRecord.StatusColumn()
		whereClause = whereClause + " AND " + statusColumn + " <> ?"
		args = append(args, Deleted)
	}

	return ds.DbMap.Select(results, whereClause, args...)
}

func getCurrentTime() time.Time {
	// Some of the SQL drivers don't store timezone information,
	// so we have to use UTC to keep the stored values equal.
	utc, _ := time.LoadLocation("UTC")
	now := time.Now().In(utc)
	return now
}

func getDialectForDriver(driverName string) gorp.Dialect {
	var gorpDialect gorp.Dialect
	switch driverName {
	case "sqlite3":
		gorpDialect = gorp.SqliteDialect{}
	default:
		gorpDialect = gorp.MySQLDialect{}
	}

	return gorpDialect
}
