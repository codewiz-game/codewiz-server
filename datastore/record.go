package datastore

import (
	"github.com/go-gorp/gorp"
	"time"
)

type LogicallyDeletable interface {
	StatusRecorder
	Keys() []interface{}
	StatusColumn() string
}

type StatusRecorder interface {
	Status() StatusCode
	SetStatus(StatusCode)
}

type CreationTimeRecorder interface {
	CreationTime() time.Time
	SetCreationTime(time.Time)
}

type DeletionTimeRecorder interface {
	DeletionTime() time.Time
	SetDeletionTime(time.Time)
}

type LastUpdateTimeRecorder interface {
	LastUpdatedTime() time.Time
	SetLastUpdatedTime(time.Time)
}

type BaseFields struct {
	ID				int64		  `db:"ID, autoincrement, primarykey"`
	Status       	StatusCode    `db:"Status"`
	CreationTime 	gorp.NullTime `db:"CreationTime"`
	LastUpdatedTime gorp.NullTime `db:"LastModified"`
	DeletionTime 	gorp.NullTime `db:"DeletionTime"`
}

type BaseRecord struct {
	BaseFields
}

func NewRecord() *BaseRecord {
	return &BaseRecord{}
}

func (record *BaseRecord) Keys() []interface{} {
	return []interface{}{ record.BaseFields.ID }
}

func (record *BaseRecord) Status() StatusCode {
	return record.BaseFields.Status
}

func (record *BaseRecord) SetStatus(status StatusCode) {
	record.BaseFields.Status = status
}

func (record *BaseRecord) StatusColumn() string {
	return "Status"
}

func (record *BaseRecord) CreationTime() time.Time {
	if !record.BaseFields.CreationTime.Valid {
		return time.Time{}
	}

	return record.BaseFields.CreationTime.Time
}

func (record *BaseRecord) SetCreationTime(time time.Time) {
	record.BaseFields.CreationTime = gorp.NullTime{Time : time, Valid : !time.IsZero()}
}

func (record *BaseRecord) DeletionTime() time.Time {
	if !record.BaseFields.DeletionTime.Valid {
		return time.Time{}
	}

	return record.BaseFields.DeletionTime.Time
}

func (record *BaseRecord) SetDeletionTime(time time.Time) {
	record.BaseFields.DeletionTime = gorp.NullTime{Time : time, Valid : !time.IsZero()}
}

func (record *BaseRecord) LastUpdatedTime() time.Time {
	if !record.BaseFields.LastUpdatedTime.Valid {
		return time.Time{}
	}

	return record.BaseFields.LastUpdatedTime.Time
}

func (record *BaseRecord) SetLastUpdatedTime(time time.Time) {
	record.BaseFields.LastUpdatedTime = gorp.NullTime{Time : time, Valid : !time.IsZero()}
}