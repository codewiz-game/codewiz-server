package datastore

import (
	"github.com/go-gorp/gorp"
	"reflect"
	"strings"
)

const (
	INVALID_ID uint64 = 0
)

var (
	activeValue  reflect.Value = reflect.ValueOf(Active)
	deletedValue reflect.Value = reflect.ValueOf(Deleted)
)

type fieldMap struct {
	Field  string
	Column string
}

type metaMap struct {
	Status   fieldMap
	Created  fieldMap
	Modified fieldMap
	Deleted  fieldMap
	Keys     []fieldMap // The gorp package offers no public methods of fields to read what the primary key of a column is, so we need to store our own list.
	Valid    bool
}

type BasicRecord struct {
	ID           uint64        `db:"ID, primarykey, autoincrement"`
	Status       StatusCode    `db:"Status" db-meta:"status"`
	CreationTime gorp.NullTime `db:"CreationTime" db-meta:"created"`
	LastModified gorp.NullTime `db:"LastModified" db-meta:"modified"`
	DeletionTime gorp.NullTime `db:"DeletionTime" db-meta:"deleted"`
}

func NewRecord() *BasicRecord {
	record := &BasicRecord{ID: INVALID_ID}
	return record
}

/**
func (record BasicRecord) ID() sql.NullInt64 {
	return sql.NullInt64{Int64: record.ID, Valid: record.ID != -1}
}
*/

// This function extracts a field from a struct given a dot-seperated path.
// For example, the path A.B.C will return the field C, which is a field of struct B,
// which is a field of struct A, which is a field of the struct represented in the other parameter.
func getFieldFromValue(fieldPath string, structVal reflect.Value) reflect.Value {

	// If no path is given, the field doesn't exist, so return a zero value
	if fieldPath == "" {
		return reflect.Value{}
	}

	fieldLevels := strings.Split(fieldPath, ".")
	for i := 0; i < len(fieldLevels); i++ {
		structVal = structVal.FieldByName(fieldLevels[i])
	}

	return structVal
}

func getMetaMap(datastore *SQLDataStore, value reflect.Value, nested bool) metaMap {

	structType := value.Type()

	// If the meta-map has already been created, return the existing one
	if existing := datastore.meta[structType]; existing.Valid {
		return existing
	}

	meta := metaMap{Valid: true}
	for i := 0; i < value.NumField(); i++ {
		currentField := structType.Field(i)
		currentFieldValue := value.FieldByName(currentField.Name)
		if currentField.Anonymous {
			// Recursively parse all of the tags in the nested structs
			nestedMeta := getMetaMap(datastore, currentFieldValue, true)
			nestedMeta.copyTo(&meta)
		} else {
			prefix := ""
			if nested {
				prefix = structType.Name() + "."
			}

			fieldPath := prefix + currentField.Name

			// Extract the primary keys from the db tag
			dbColumn := ""
			if dbTag := currentField.Tag.Get("db"); dbTag != "" {

				dbTagAttrs := strings.Split(dbTag, ",")
				if len(dbTagAttrs) > 0 {
					dbColumn = dbTagAttrs[0]
					for j := 1; j < len(dbTagAttrs); j++ { // ignore the first flag, since this is always the field name
						trimmedFlag := strings.TrimSpace(dbTagAttrs[j])
						if trimmedFlag == "primarykey" {
							meta.Keys = append(meta.Keys, fieldMap{Field: fieldPath, Column: dbColumn})
						}
					}
				}
			}

			// Extract all of the other fields from the db-meta tag
			if dbMetaTag := currentField.Tag.Get("db-meta"); dbMetaTag != "" {
				trimmedTag := strings.TrimSpace(dbMetaTag)
				switch trimmedTag {
				case "status":
					meta.Status = fieldMap{Field: fieldPath, Column: dbColumn}
				case "created":
					meta.Created = fieldMap{Field: fieldPath, Column: dbColumn}
				case "modified":
					meta.Modified = fieldMap{Field: fieldPath, Column: dbColumn}
				case "deleted":
					meta.Deleted = fieldMap{Field: fieldPath, Column: dbColumn}
				}
			}
		}
	}

	datastore.meta[structType] = meta
	return meta
}

func (from *metaMap) copyTo(to *metaMap) {
	// Only copy a field if it is not already set - this ensures values set by nested structs can be overridden

	if to.Status.Field == "" && from.Status.Field != "" {
		to.Status = from.Status
	}

	if to.Created.Field == "" && from.Created.Field != "" {
		to.Created = from.Created
	}

	if to.Modified.Field == "" && from.Modified.Field != "" {
		to.Modified = from.Modified
	}

	if to.Deleted.Field == "" && from.Deleted.Field != "" {
		to.Deleted = from.Deleted
	}

	for _, key := range from.Keys {
		to.Keys = append(to.Keys, key)
	}
}
