package datastore

import (
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
	BaseRecord
	String  string `db:"StringField"`
	Integer int    `db:"IntegerField"`
}

func (result *testResult) LastInsertId() (int64, error) {
	return result.lastInsertId, result.err
}

func (result *testResult) RowsAffected() (int64, error) {
	return result.rowsAffected, result.err
}


func TestDBStore_Insert_SucceedsOnValidRecord(t *testing.T) {
	ds, err := initTestDataStore()
	if err != nil {
		t.Fatal(err)
	}

	// Create a valid record
	record := &testRecord{BaseRecord: *NewRecord(), String: "ABC", Integer: 20}

	// Insert the record using the datastore
	preInsertTime := time.Now()
	err = ds.Insert(record)
	if err != nil {
		t.Fatal(err)
	}

	// Ensure that the creation time was updated locally
	if record.CreationTime().IsZero() || record.CreationTime().Before(preInsertTime) {
		t.Fatalf("Creation time did not update")
	}

	// Ensure that the last updated time was updated locally
	if record.LastUpdatedTime().IsZero() || record.LastUpdatedTime().Before(preInsertTime) {
		t.Fatalf("Last modifed time did not update")
	}

	// Select the record that was just inserted from the DB,
	// and then assert that all of the fields inserted correctly.
	inserted := &testRecord{}
	err = ds.DbMap.SelectOne(inserted, "SELECT * FROM Test") // no need to specify ID, as there should only be one record
	if err != nil {
		t.Fatal(err)
	}

	assertEquals(record, inserted, t)
	closeTestDatastore(ds)
}

func TestDBStore_Insert_FailsOnInvalidRecord(t *testing.T) {
	ds, err := initTestDataStore()
	if err != nil {
		t.Fatal(err)
	}

	// Create a record that has an IntegerField that breaks the constraint
	record := &testRecord{BaseRecord: *NewRecord(), String: "ABC", Integer: -1}

	err = ds.Insert(record)
	if err == nil { // expect this operation to fail
		t.Fatalf("Insert operation was expected to return an error")
	}

	// Ensure that the creation and last modified times are still invalid
	if !record.CreationTime().IsZero() || !record.LastUpdatedTime().IsZero() {
		t.Fatalf("Creation and/or last modified times did not roll back")
	}

	closeTestDatastore(ds)
}

func TestDBStore_Insert_FailsOnSameRecord(t *testing.T) {
	ds, err := initTestDataStore()
	if err != nil {
		t.Fatal(err)
	}

	// Create a valid Record
	record := &testRecord{BaseRecord: *NewRecord(), String: "ABC", Integer: 20}

	// Insert the record into the database
	err = ds.Insert(record)
	if err != nil {
		t.Fatal(err)
	}

	previousLastModified := record.LastUpdatedTime()
	previousCreationTime := record.CreationTime()

	// Attempt to re-insert the exact same record
	err = ds.Insert(record)
	if err == nil { // expect this operation to fail
		t.Fatalf("Did not expect insertion operation to succeed")
	}

	// Ensure that the creation and last modifed times on the local copy are the same as they were after the successful insert
	if (record.LastUpdatedTime() != previousLastModified) || (record.CreationTime() != previousCreationTime) {
		t.Fatalf("Did not expect last modified and/or creation time to update")
	}

	// Ensure that there is only one persisted copy, and that its fields are equal to the local copy
	inserted := &testRecord{}
	err = ds.DbMap.SelectOne(inserted, "SELECT * FROM Test") // no need to specify ID, as there should only be one record
	if err != nil {
		t.Fatal(err)
	}

	assertEquals(record, inserted, t)

	closeTestDatastore(ds)
}

func TestDBStore_Insert_FailsOnExistingPrimaryKey(t *testing.T) {
	ds, err := initTestDataStore()
	if err != nil {
		t.Fatal(err)
	}

	// Create a valid Record
	record := &testRecord{BaseRecord: *NewRecord(), String: "ABC", Integer: 20}

	// Insert the record into the database
	err = ds.Insert(record)
	if err != nil {
		t.Fatal(err)
	}

	updated := &testRecord{BaseRecord: *NewRecord(), String: "DEF", Integer: 35}
	updated.ID = record.ID

	// Attempt to insert the 'updated' record that has the same primary key
	err = ds.Insert(updated)
	if err == nil { // expect this operation to fail
		t.Fatalf("Did not expect insertion operation to succeed")
	}

	// Ensure that the creation and last modifed times on the updated copy are still not set (ie. invalid)
	if !(updated.CreationTime().IsZero() && updated.LastUpdatedTime().IsZero()) {
		t.Fatalf("Did not expect creation time and/or last modified time to be set")
	}

	// Ensure that there is only one persisted copy, and that its fields are equal to the local copy
	inserted := &testRecord{}
	err = ds.DbMap.SelectOne(inserted, "SELECT * FROM Test") // no need to specify ID, as there should only be one record
	if err != nil {
		t.Fatal(err)
	}

	assertEquals(record, inserted, t)

	closeTestDatastore(ds)
}

func TestDBStore_Insert_PreviouslyDeletedRecord(t *testing.T) {
	ds, err := initTestDataStore()
	if err != nil {
		t.Fatal(err)
	}

	// Insert a record marked as Deleted
	record := &testRecord{BaseRecord: *NewRecord(), String: "ABC", Integer: 20}
	record.SetStatus(Deleted)
	ds.DbMap.Insert(record)

	// Attempt to insert the same record
	err = ds.Insert(record)

	// Ensure that the update succeeded
	if err != nil {
		t.Fatal(err)
	}

	// Ensure that there is only one persisted copy, and that its fields are equal to the local copy
	inserted := &testRecord{}
	err = ds.DbMap.SelectOne(inserted, "SELECT * FROM Test") // no need to specify ID, as there should only be one record
	if err != nil {
		t.Fatal(err)
	}

	assertEquals(record, inserted, t)

	closeTestDatastore(ds)
}

func TestDBStore_Update_SucceedsOnValidRecord(t *testing.T) {
	ds, err := initTestDataStore()
	if err != nil {
		t.Fatal(err)
	}

	// Create a valid record
	record := &testRecord{BaseRecord: *NewRecord(), String: "ABC", Integer: 20}

	// Insert the record into the database
	err = ds.DbMap.Insert(record)
	if err != nil {
		t.Fatal(err)
	}

	// Modify the record and request an update
	preUpdateTime := time.Now()
	record.String = "DEF"
	record.Integer = 45
	err = ds.Update(record)
	if err != nil {
		t.Fatal(err)
	}

	// Ensure that the last modified time was updated locally
	if record.LastUpdatedTime().Before(preUpdateTime) {
		t.Fatalf("Did not expect last modified time to be updated")
	}

	// Select the record from the DB and assert that all of fields updated correctly
	updated := &testRecord{}
	err = ds.DbMap.SelectOne(updated, "SELECT * FROM Test") // no need to specify ID, as there should only be one record
	if err != nil {
		t.Fatal(err)
	}

	assertEquals(*record, *updated, t)

	closeTestDatastore(ds)
}

func TestDBStore_Update_FailsOnInvalidRecord(t *testing.T) {
	ds, err := initTestDataStore()
	if err != nil {
		t.Fatal(err)
	}

	// Create a valid record
	record := &testRecord{BaseRecord: *NewRecord(), String: "ABC", Integer: 20}

	// Insert the record into the database
	err = ds.DbMap.Insert(record)
	if err != nil {
		t.Fatal(err)
	}

	// Modify the record to break the IntegerField constraint and request an update
	previousUpdateTime := record.LastUpdatedTime()
	record.String = "DEF"
	record.Integer = -1
	err = ds.Update(record)
	if err == nil { // expect this operation to fail
		t.Fatalf("Did not expect update operation to succeed")
	}

	// Ensure that the last update time is equal to what it was before the update (ie. the insertion time)
	if !(record.LastUpdatedTime() == previousUpdateTime && record.LastUpdatedTime() == record.CreationTime()) {
		t.Fatalf("Did not expect last modified time to be updated")
	}

	closeTestDatastore(ds)
}

func TestDBStore_Update_FailsOnMissingRecord(t *testing.T) {
	ds, err := initTestDataStore()
	if err != nil {
		t.Fatal(err)
	}

	// Create a valid record
	record := &testRecord{BaseRecord: *NewRecord(), String: "ABC", Integer: 20}

	// Attempt to update the record in the datastore
	record.String = "DEF"
	record.Integer = 45
	err = ds.Update(record)
	if err == nil { // expect this operation to fail
		t.Fatalf("Did not expect update operation to succeed")
	}

	// Ensure that the last modified time is still not set (ie. invalid)
	if !record.LastUpdatedTime().IsZero() {
		t.Fatalf("Did not expect last modified time to be set")
	}

	closeTestDatastore(ds)
}

func TestDBStore_Delete_SucceedsOnExistingRecord(t *testing.T) {
	ds, err := initTestDataStore()
	if err != nil {
		t.Fatal(err)
	}

	// Create a valid record
	record := &testRecord{BaseRecord: *NewRecord(), String: "ABC", Integer: 20}

	// Insert the record into the database
	err = ds.DbMap.Insert(record)
	if err != nil {
		t.Fatal(err)
	}

	// Delete the same record
	preDeletionTime := time.Now()
	err = ds.Delete(record)
	if err != nil {
		t.Fatal(err)
	}

	// Ensure that the deletion time updated
	if record.DeletionTime().IsZero() || record.DeletionTime().Before(preDeletionTime) {
		t.Fatalf("Deletion time was not updated!")
	}

	// Ensure that the status of the record is now set to Deleted
	if record.Status() != Deleted {
		t.Fatalf("Status was not marked as deleted")
	}

	// Ensure that a search for deleted records returns this record
	deleted := &testRecord{}
	err = ds.DbMap.SelectOne(deleted, "SELECT * FROM Test WHERE Status = ?", Deleted)
	if err != nil {
		t.Fatal(err)
	}

	assertEquals(*record, *deleted, t)

	closeTestDatastore(ds)
}

func TestDBStore_Delete_FailsOnMissingRecord(t *testing.T) {
	ds, err := initTestDataStore()
	if err != nil {
		t.Fatal(err)
	}

	// Create a valid record
	record := &testRecord{BaseRecord: *NewRecord(), String: "ABC", Integer: 20}

	// Request that the record be deleted from the datastore
	err = ds.Delete(record)
	if err == nil { // expect this operation to fail
		t.Fatalf("Deletion operation was expected to fail")
	}

	// Ensure that the deletion time is not set (invalid)
	if !record.DeletionTime().IsZero() {
		t.Fatalf("Did not expect deletion time to be set")
	}

	// Ensure that the record is still in a Transient state
	if record.Status() != Transient {
		t.Fatalf("Did not expect status field to be updated")
	}

	closeTestDatastore(ds)
}

func TestDBStore_Get_IgnoresLogicallyDeletedRecords(t *testing.T) {
	ds, err := initTestDataStore()
	if err != nil {
		t.Fatal(err)
	}

	// Insert one record that is marked as Deleted, and another that is Active
	activeRecord := &testRecord{BaseRecord: *NewRecord(), String: "ABC", Integer: 20}
	activeRecord.SetStatus(Active)
	ds.DbMap.Insert(activeRecord)

	deletedRecord := &testRecord{BaseRecord: *NewRecord(), String: "DEF", Integer: 45}
	deletedRecord.SetStatus(Deleted)
	ds.DbMap.Insert(deletedRecord)

	// Attempt to get the records from the datastore
	result, err := ds.Get(testRecord{}, "SELECT * FROM Test WHERE StringField = ?", "ABC")
	if err != nil || result == nil {
		t.Fatalf("No results returned when searching for an Active record")
	}

	insertedRecord := result.(*testRecord)
	assertEquals(*activeRecord, *insertedRecord, t)

	insertedDeleted, err := ds.Get(testRecord{}, "SELECT * FROM Test WHERE StringField = ?", "DEF")
	if err != nil {
		t.Fatalf("Unexpected error when trying to query for a Deleted record")
	}

	if insertedDeleted != nil {
		t.Fatalf("No result should be returned when querying for a deleted record")
	}

	closeTestDatastore(ds)
}

func TestDBStore_Get_MoreThanOneResult(t *testing.T) {
	ds, err := initTestDataStore()
	if err != nil {
		t.Fatal(err)
	}

	// Insert two Active records
	firstRecord := &testRecord{BaseRecord: *NewRecord(), String: "ABC", Integer: 20}
	firstRecord.SetStatus(Active)
	ds.DbMap.Insert(firstRecord)

	secondRecord := &testRecord{BaseRecord: *NewRecord(), String: "DEF", Integer: 45}
	secondRecord.SetStatus(Active)
	ds.DbMap.Insert(secondRecord)

	result, err := ds.Get(testRecord{}, "SELECT * FROM Test")
	if result != nil || err == nil {
		t.Fatalf("Expected an error to be thrown with no returned result")
	}

	closeTestDatastore(ds)
}


func initTestDataStore() (*DBStore, error) {
	ds, err := NewDBStore("sqlite3", ":memory:")
	if err != nil {
		return nil, err
	}

	_, err = ds.Exec(`
		CREATE TABLE Test (
		ID INTEGER PRIMARY KEY,
		Status INTEGER NOT NULL,
		CreationTime DATETIME,
		DeletionTime DATETIME,
		LastModified DATETIME, 
		StringField VARCHAR(255),
		IntegerField INTEGER
		CHECK (IntegerField > 0)
	);`)

	if err != nil {
		return nil, err
	}
	
	ds.AddTableWithName(testRecord{}, "Test")
	return ds, nil
}

func closeTestDatastore(ds *DBStore) {
	ds.DbMap.Db.Close()
}

func assertEquals(a interface{}, b interface{}, t *testing.T) {
	if !reflect.DeepEqual(a, b) {
		t.Fatalf("Persisted copy [%v] is not equal to local copy [%v]", a, b)
	}
}
