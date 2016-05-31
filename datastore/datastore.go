package datastore

import (
	"database/sql"
	"errors"
	"github.com/go-gorp/gorp"
	"reflect"
	"strings"
	"time"
)

type StatusCode uint
type Dialect uint

const (
	Transient StatusCode = iota // not yet inserted - the zero-value for StatusCode
	Active
	Deleted

	UnsupportedDialect Dialect = iota
	MySQLDialect
	SqliteDialect
)

type SQLDataStore struct {
	*gorp.DbMap
	meta map[reflect.Type]metaMap
}

func NewDatastore(db *sql.DB, driverName string) *SQLDataStore {
	gorpDialect := getDialectForDriver(driverName)
	dbMap := &gorp.DbMap{Db: db, Dialect : gorpDialect}
	return &SQLDataStore{DbMap: dbMap, meta: make(map[reflect.Type]metaMap)}
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

func (ds *SQLDataStore) AddTableWithName(model interface{}, name string) {

	structVal := reflect.ValueOf(model)
	structType := structVal.Type()

	// Load the necessary field/column information into our map
	ds.meta[structType] = getMetaMap(ds, structVal, false)

	// Pass the struct to gorp for it to register in its own maps
	ds.DbMap.AddTableWithName(model, name)
}

func (ds *SQLDataStore) Insert(record interface{}) error {

	structVal := reflect.ValueOf(record).Elem()
	structType := structVal.Type()
	meta := ds.meta[structType]

	// gorp does not include the auto-increment primary keys in it's insertions, which means
	// that the same record can be passed in multiple times and will generate
	// additional rows. To fix this, we need to query for an matching record first and assert
	// that no Active records are found.
	existing, err := ds.getPersistedRecord(record)
	if err != nil {
		return err
	}

	if existing != nil {
		existingVal := reflect.ValueOf(existing).Elem()
		existingStatusField := getFieldFromValue(meta.Status.Field, existingVal)
		if existingStatusField.IsValid() {
			if status := existingStatusField.Interface(); status == Active {
				// If the existing record is still active, throw an error
				return errors.New("A record with the same primary key already existing in the data store")
			} else if status == Deleted {
				// If the existing record is Deleted, change it back to Active
				existingStatusField.Set(activeValue)
				err = ds.Update(record)
				return err
			}
		} else {
			// If the status field isn't being used, the record in the database is active
			return errors.New("A record with the same primary key already existing in the data store")
		}
	}

	statusField := getFieldFromValue(meta.Status.Field, structVal)
	creationField := getFieldFromValue(meta.Created.Field, structVal)
	modifiedField := getFieldFromValue(meta.Modified.Field, structVal)

	// Update the status, creation time and last modifed before insertion, so that the changes can be persisted
	var previousStatus interface{}
	if statusField.IsValid() {
		previousStatus = statusField.Interface()
		statusField.Set(activeValue)
	}

	now := getCurrentTimeValue()
	var previousCreationTime interface{}
	if creationField.IsValid() {
		previousCreationTime = creationField.Interface()
		creationField.Set(now)
	}

	var previousLastModified interface{}
	if modifiedField.IsValid() {
		previousLastModified = modifiedField.Interface()
		modifiedField.Set(now)
	}

	err = ds.DbMap.Insert(record)

	if err != nil {
		// If the insertion failed, revert the changes made to
		// the status, creation and last modified time to keep the data accurate
		if statusField.IsValid() {
			statusField.Set(reflect.ValueOf(previousStatus))
		}
		if creationField.IsValid() {
			creationField.Set(reflect.ValueOf(previousCreationTime))
		}
		if modifiedField.IsValid() {
			modifiedField.Set(reflect.ValueOf(previousLastModified))
		}

		return err
	}

	return nil
}

func (ds *SQLDataStore) Update(record interface{}) error {

	structVal := reflect.ValueOf(record).Elem()
	structType := structVal.Type()
	meta := ds.meta[structType]

	modifiedField := getFieldFromValue(meta.Modified.Field, structVal)

	now := getCurrentTimeValue()
	var previousLastModified interface{}
	if modifiedField.IsValid() {
		previousLastModified = modifiedField.Interface()
		modifiedField.Set(now)
	}

	count, err := ds.DbMap.Update(record)
	if err == nil && count == 0 {
		err = errors.New("No records were affected by the update operation.")
	}

	if err != nil {
		if modifiedField.IsValid() {
			modifiedField.Set(reflect.ValueOf(previousLastModified))
		}
		return err
	}

	return nil
}

func (ds *SQLDataStore) Delete(record interface{}) error {

	structVal := reflect.ValueOf(record).Elem()
	structType := structVal.Type()
	meta := ds.meta[structType]

	deletionField := getFieldFromValue(meta.Deleted.Field, structVal)
	statusField := getFieldFromValue(meta.Status.Field, structVal)

	now := getCurrentTimeValue()
	var previousDeletionTime interface{}
	if deletionField.IsValid() {
		previousDeletionTime = deletionField.Interface()
		deletionField.Set(now)
	}

	var previousStatus interface{}
	if statusField.IsValid() {
		previousStatus = statusField.Interface()
		statusField.Set(deletedValue)
	}
	err := ds.Update(record)
	if err != nil {
		if deletionField.IsValid() {
			deletionField.Set(reflect.ValueOf(previousDeletionTime))
		}

		if statusField.IsValid() {
			statusField.Set(reflect.ValueOf(previousStatus))
		}
		return err
	}

	return nil
}

func (ds *SQLDataStore) Get(model interface{}, whereClause string, args ...interface{}) (interface{}, error) {

	results, err := ds.Select(model, whereClause, args...)
	if err != nil {
		return nil, err
	}

	resultsVal := reflect.ValueOf(results)
	numResults := resultsVal.Len()

	if numResults == 0 {

		return nil, nil
	}

	if numResults != 1 {
		return nil, errors.New("More than one matching record was found for the given query")
	}

	record := resultsVal.Index(0).Interface()

	return record, nil
}

func (ds *SQLDataStore) Select(model interface{}, whereClause string, args ...interface{}) ([]interface{}, error) {

	structType := reflect.TypeOf(model)
	for structType.Kind() == reflect.Ptr {
		structType = structType.Elem()
	}

	meta := ds.meta[structType]

	statusColumn := meta.Status.Column
	if statusColumn != "" {
		whereClause = whereClause + " AND " + statusColumn + " <> ?"
		args = append(args, Deleted)
	}

	model = reflect.New(structType).Interface()
	results, err := ds.DbMap.Select(reflect.New(structType).Interface(), whereClause, args...)
	if err != nil {
		return nil, err
	}

	return results, nil
}

func getCurrentTimeValue() reflect.Value {
	// Some of the SQL drivers don't store timezone information,
	// so we have to use UTC to keep the stored values equal.
	utc, _ := time.LoadLocation("UTC")
	now := time.Now().In(utc)
	return reflect.ValueOf(gorp.NullTime{Time: now, Valid: true})
}

// TODO: this is extremely similar to the existing Get implementation, but without checking the Status column
// The common code should eventually be moved into a seperate function that they both call
func (ds *SQLDataStore) getPersistedRecord(record interface{}) (interface{}, error) {
	structVal := reflect.ValueOf(record).Elem()
	structType := structVal.Type()

	table, err := ds.DbMap.TableFor(structType, true)
	if err != nil {
		return nil, err
	}

	keys := ds.meta[structType].Keys
	setters := make([]string, len(keys))
	args := make([]interface{}, len(keys))
	for index, key := range keys {
		setters[index] = key.Column + " = ?"
		keyField := getFieldFromValue(key.Field, structVal)
		args[index] = keyField.Interface()
	}

	whereClause := "SELECT * FROM " + table.TableName + " WHERE " + strings.Join(setters, " AND ")
	model := reflect.Zero(structType).Interface()
	results, err := ds.DbMap.Select(model, whereClause, args...)
	if err != nil {
		return nil, err
	}

	resultsVal := reflect.ValueOf(results)

	if resultsVal.Len() == 0 {
		return nil, nil
	}

	return resultsVal.Index(0).Interface(), nil
}
