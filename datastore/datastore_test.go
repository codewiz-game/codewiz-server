package datastore

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"reflect"
	"testing"
	"time"
)

type testResult struct {
	lastInsertId int64
	rowsAffected int64
	err          error
}

type testRecord struct {
	BasicRecord
	String  string `db:"StringField"`
	Integer int    `db:"IntegerField"`
}

func (result *testResult) LastInsertId() (int64, error) {
	return result.lastInsertId, result.err
}

func (result *testResult) RowsAffected() (int64, error) {
	return result.rowsAffected, result.err
}

func initTestDataStore() (*SQLDataStore, error) {
	db, err := sql.Open("sqlite3", ":memory:")

	sql := `CREATE TABLE Test (
		ID INTEGER PRIMARY KEY,
		Status INTEGER NOT NULL,
		CreationTime DATETIME,
		DeletionTime DATETIME,
		LastModified DATETIME, 
		StringField VARCHAR(255),
		IntegerField INTEGER
		CHECK (IntegerField > 0)
	);`

	_, err = db.Exec(sql)
	if err != nil {
		return nil, err
	}

	ds := NewDataStore(db, "sqlite3")
	ds.AddTableWithName(testRecord{}, "Test")
	return ds, nil
}

func closeTestDatastore(ds *SQLDataStore) {
	ds.DbMap.Db.Close()
}

func TestSQLDataStore_Insert_SucceedsOnValidRecord(t *testing.T) {
	ds, err := initTestDataStore()
	if err != nil {
		t.Error(err)
	}

	// Create a valid record
	record := &testRecord{BasicRecord: *NewRecord(), String: "ABC", Integer: 20}

	// Insert the record using the datastore
	preInsertTime := time.Now()
	err = ds.Insert(record)
	if err != nil {
		t.Error(err)
	}

	// Ensure that the creation time was updated locally
	if !record.CreationTime.Valid || record.CreationTime.Time.Before(preInsertTime) {
		t.Errorf("Creation time did not update")
	}

	// Ensure that the last modified time was updated locally
	if !record.LastModified.Valid || record.LastModified.Time.Before(preInsertTime) {
		t.Errorf("Last modifed time did not update")
	}

	// Ensure that the ID was updated
	if record.ID == INVALID_ID {
		t.Errorf("ID did not update")
	}

	// Select the record that was just inserted from the DB,
	// and then assert that all of the fields inserted correctly.
	inserted := &testRecord{}
	err = ds.DbMap.SelectOne(inserted, "SELECT * FROM Test") // no need to specify ID, as there should only be one record
	if err != nil {
		t.Error(err)
	}

	assertEquals(record, inserted, t)

	closeTestDatastore(ds)
}

func TestSQLDataStore_Insert_FailsOnInvalidRecord(t *testing.T) {
	ds, err := initTestDataStore()
	if err != nil {
		t.Error(err)
	}

	// Create a record that has an IntegerField that breaks the constraint
	record := &testRecord{BasicRecord: *NewRecord(), String: "ABC", Integer: -1}

	err = ds.Insert(record)
	if err == nil { // expect this operation to fail
		t.Errorf("Insert operation was expected to return an error")
	}

	// Ensure that the creation and last modified times are still invalid
	if record.CreationTime.Valid || record.LastModified.Valid {
		t.Errorf("Creation and/or last modified times did not roll back")
	}

	closeTestDatastore(ds)
}

func TestSQLDataStore_Insert_FailsOnSameRecord(t *testing.T) {
	ds, err := initTestDataStore()
	if err != nil {
		t.Error(err)
	}

	// Create a valid Record
	record := &testRecord{BasicRecord: *NewRecord(), String: "ABC", Integer: 20}

	// Insert the record into the database
	err = ds.Insert(record)
	if err != nil {
		t.Error(err)
	}

	previousLastModified := record.LastModified
	previousCreationTime := record.CreationTime

	// Attempt to re-insert the exact same record
	err = ds.Insert(record)
	if err == nil { // expect this operation to fail
		t.Errorf("Did not expect insertion operation to succeed")
	}

	// Ensure that the creation and last modifed times on the local copy are the same as they were after the successful insert
	if !(record.LastModified.Valid && record.LastModified == previousLastModified) ||
		!(record.CreationTime.Valid && record.CreationTime == previousCreationTime) {
		t.Errorf("Did not expect last modified and/or creation time to update")
	}

	// Ensure that there is only one persisted copy, and that its fields are equal to the local copy
	inserted := &testRecord{}
	err = ds.DbMap.SelectOne(inserted, "SELECT * FROM Test") // no need to specify ID, as there should only be one record
	if err != nil {
		t.Error(err)
	}

	assertEquals(record, inserted, t)

	closeTestDatastore(ds)
}

func TestSQLDataStore_Insert_FailsOnExistingPrimaryKey(t *testing.T) {
	ds, err := initTestDataStore()
	if err != nil {
		t.Error(err)
	}

	// Create a valid Record
	record := &testRecord{BasicRecord: *NewRecord(), String: "ABC", Integer: 20}

	// Insert the record into the database
	err = ds.Insert(record)
	if err != nil {
		t.Error(err)
	}

	updated := &testRecord{BasicRecord: *NewRecord(), String: "DEF", Integer: 35}
	updated.ID = record.ID

	// Attempt to insert the 'updated' record that has the same primary key
	err = ds.Insert(updated)
	if err == nil { // expect this operation to fail
		t.Errorf("Did not expect insertion operation to succeed")
	}

	// Ensure that the creation and last modifed times on the updated copy are still not set (ie. invalid)
	if updated.CreationTime.Valid || updated.LastModified.Valid {
		t.Errorf("Did not expect creation time and/or last modified time to be set")
	}

	// Ensure that there is only one persisted copy, and that its fields are equal to the local copy
	inserted := &testRecord{}
	err = ds.DbMap.SelectOne(inserted, "SELECT * FROM Test") // no need to specify ID, as there should only be one record
	if err != nil {
		t.Error(err)
	}

	assertEquals(record, inserted, t)

	closeTestDatastore(ds)
}

func TestSQLDataStore_Insert_PreviouslyDeletedRecord(t *testing.T) {
	ds, err := initTestDataStore()
	if err != nil {
		t.Error(err)
	}

	// Insert a record marked as Deleted
	record := &testRecord{BasicRecord: *NewRecord(), String: "ABC", Integer: 20}
	record.Status = Deleted
	ds.DbMap.Insert(record)

	// Attempt to insert the same record
	err = ds.Insert(record)

	// Ensure that the update succeeded
	if err != nil {
		t.Error(err)
	}

	// Ensure that there is only one persisted copy, and that its fields are equal to the local copy
	inserted := &testRecord{}
	err = ds.DbMap.SelectOne(inserted, "SELECT * FROM Test") // no need to specify ID, as there should only be one record
	if err != nil {
		t.Error(err)
	}

	assertEquals(record, inserted, t)

	closeTestDatastore(ds)
}

func TestSQLDataStore_Update_SucceedsOnValidRecord(t *testing.T) {
	ds, err := initTestDataStore()
	if err != nil {
		t.Error(err)
	}

	// Create a valid record
	record := &testRecord{BasicRecord: *NewRecord(), String: "ABC", Integer: 20}

	// Insert the record into the database
	err = ds.DbMap.Insert(record)
	if err != nil {
		t.Error(err)
	}

	// Modify the record and request an update
	preUpdateTime := time.Now()
	record.String = "DEF"
	record.Integer = 45
	err = ds.Update(record)
	if err != nil {
		t.Error(err)
	}

	// Ensure that the last modified time was updated locally
	if !record.LastModified.Valid || record.LastModified.Time.Before(preUpdateTime) {
		t.Errorf("Did not expect last modified time to be updated")
	}

	// Select the record from the DB and assert that all of fields updated correctly
	updated := &testRecord{}
	err = ds.DbMap.SelectOne(updated, "SELECT * FROM Test") // no need to specify ID, as there should only be one record
	if err != nil {
		t.Error(err)
	}

	assertEquals(*record, *updated, t)

	closeTestDatastore(ds)
}

func TestSQLDataStore_Update_FailsOnInvalidRecord(t *testing.T) {
	ds, err := initTestDataStore()
	if err != nil {
		t.Error(err)
	}

	// Create a valid record
	record := &testRecord{BasicRecord: *NewRecord(), String: "ABC", Integer: 20}

	// Insert the record into the database
	err = ds.DbMap.Insert(record)
	if err != nil {
		t.Error(err)
	}

	// Modify the record to break the IntegerField constraint and request an update
	previousUpdateTime := record.LastModified
	record.String = "DEF"
	record.Integer = -1
	err = ds.Update(record)
	if err == nil { // expect this operation to fail
		t.Errorf("Did not expect update operation to succeed")
	}

	// Ensure that the last update time is equal to what it was before the update (ie. the insertion time)
	if !(record.LastModified == previousUpdateTime && record.LastModified == record.CreationTime) {
		t.Errorf("Did not expect last modified time to be updated")
	}

	closeTestDatastore(ds)
}

func TestSQLDataStore_Update_FailsOnMissingRecord(t *testing.T) {
	ds, err := initTestDataStore()
	if err != nil {
		t.Error(err)
	}

	// Create a valid record
	record := &testRecord{BasicRecord: *NewRecord(), String: "ABC", Integer: 20}

	// Attempt to update the record in the datastore
	record.String = "DEF"
	record.Integer = 45
	err = ds.Update(record)
	if err == nil { // expect this operation to fail
		t.Errorf("Did not expect update operation to succeed")
	}

	// Ensure that the last modified time is still not set (ie. invalid)
	if record.LastModified.Valid {
		t.Errorf("Did not expect last modified time to be set")
	}

	closeTestDatastore(ds)
}

func TestSQLDataStore_Delete_SucceedsOnExistingRecord(t *testing.T) {
	ds, err := initTestDataStore()
	if err != nil {
		t.Error(err)
	}

	// Create a valid record
	record := &testRecord{BasicRecord: *NewRecord(), String: "ABC", Integer: 20}

	// Insert the record into the database
	err = ds.DbMap.Insert(record)
	if err != nil {
		t.Error(err)
	}

	// Delete the same record
	preDeletionTime := time.Now()
	err = ds.Delete(record)
	if err != nil {
		t.Error(err)
	}

	// Ensure that the deletion time updated
	if !record.DeletionTime.Valid || record.DeletionTime.Time.Before(preDeletionTime) {
		t.Errorf("Deletion time was not updated!")
	}

	// Ensure that the status of the record is now set to Deleted
	if record.Status != Deleted {
		t.Errorf("Status was not marked as deleted")
	}

	// Ensure that a search for deleted records returns this record
	deleted := &testRecord{}
	err = ds.DbMap.SelectOne(deleted, "SELECT * FROM Test WHERE Status = ?", Deleted)
	if err != nil {
		t.Error(err)
	}

	assertEquals(*record, *deleted, t)

	closeTestDatastore(ds)
}

func TestSQLDataStore_Delete_FailsOnMissingRecord(t *testing.T) {
	ds, err := initTestDataStore()
	if err != nil {
		t.Error(err)
	}

	// Create a valid record
	record := &testRecord{BasicRecord: *NewRecord(), String: "ABC", Integer: 20}

	// Request that the record be deleted from the datastore
	err = ds.Delete(record)
	if err == nil { // expect this operation to fail
		t.Errorf("Deletion operation was expected to fail")
	}

	// Ensure that the deletion time is not set (invalid)
	if record.DeletionTime.Valid {
		t.Errorf("Did not expect deletion time to be set")
	}

	// Ensure that the record is still in a Transient state
	if record.Status != Transient {
		t.Errorf("Did not expect status field to be updated")
	}

	closeTestDatastore(ds)
}

func TestSQLDataStore_Get_IgnoresLogicallyDeletedRecords(t *testing.T) {
	ds, err := initTestDataStore()
	if err != nil {
		t.Error(err)
	}

	// Insert one record that is marked as Deleted, and another that is Active
	activeRecord := &testRecord{BasicRecord: *NewRecord(), String: "ABC", Integer: 20}
	activeRecord.Status = Active
	ds.DbMap.Insert(activeRecord)

	deletedRecord := &testRecord{BasicRecord: *NewRecord(), String: "DEF", Integer: 45}
	deletedRecord.Status = Deleted
	ds.DbMap.Insert(deletedRecord)

	// Attempt to get the records from the datastore
	result, err := ds.Get(testRecord{}, "SELECT * FROM Test WHERE StringField = ?", "ABC")
	if err != nil || result == nil {
		t.Errorf("No results returned when searching for an Active record")
	}

	insertedRecord := result.(*testRecord)
	assertEquals(*activeRecord, *insertedRecord, t)

	insertedDeleted, err := ds.Get(testRecord{}, "SELECT * FROM Test WHERE StringField = ?", "DEF")
	if err != nil {
		t.Errorf("Unexpected error when trying to query for a Deleted record")
	}

	if insertedDeleted != nil {
		t.Errorf("No result should be returned when querying for a deleted record")
	}

	closeTestDatastore(ds)
}

func TestSQLDataStore_Get_MoreThanOneResult(t *testing.T) {
	ds, err := initTestDataStore()
	if err != nil {
		t.Error(err)
	}

	// Insert two Active records
	firstRecord := &testRecord{BasicRecord: *NewRecord(), String: "ABC", Integer: 20}
	firstRecord.Status = Active
	ds.DbMap.Insert(firstRecord)

	secondRecord := &testRecord{BasicRecord: *NewRecord(), String: "DEF", Integer: 45}
	secondRecord.Status = Active
	ds.DbMap.Insert(secondRecord)

	result, err := ds.Get(testRecord{}, "SELECT * FROM Test")
	if result != nil || err == nil {
		t.Errorf("Expected an error to be thrown with no returned result")
	}

	closeTestDatastore(ds)
}

func assertEquals(a interface{}, b interface{}, t *testing.T) {
	if !reflect.DeepEqual(a, b) {
		t.Errorf("Persisted copy [%v] is not equal to local copy [%v]", a, b)
	}
}
