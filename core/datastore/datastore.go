package datastore

type StatusCode int

const (
	Transient StatusCode = iota // not yet inserted - the zero-value for StatusCode
	Active
	Deleted
)

type DataStore interface {
	AddRepositoryWithName(model interface{}, name string)
	Insert(record interface{}) error
	Update(record interface{}) error
	Delete(record interface{}) error
	Get(whereClause string, args ...interface{}) (interface{}, error)
	Select(whereClause string, args ...interface{}) (interface{}, error)
}
