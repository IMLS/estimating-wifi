package interfaces

import (
	"github.com/jmoiron/sqlx"
)

type Database interface {
	// NewDB(path string) *sqlx.DB
	Open()
	Close()
	GetPtr() *sqlx.DB
	GetPath() string
	InitTable(name string) Table
	CreateTableFromStruct(s interface{}) Table
	RemoveTable(name string)
	CheckTableExists(name string) bool
	ListTables() []string
	GetTableFromStruct(s interface{}) Table
	GetTableByName(name string) Table
	Query(string) (*sqlx.Rows, error)
}

type Table interface {
	AddColumn(name string, sqltype string)
	Create()
	InsertStruct(s interface{})
	InsertMany(s []interface{})
	Drop()
	GetIntegerType() string
	GetTextType() string
	GetDateType() string
	GetDB() Database
	//FieldToDate(field string) time.Time
	//DateToField(t time.Time) string
	GetIntegerField(string) int
	GetTextField(string) string
	SetIntegerField(name string, value int)
	SetTextField(name string, value string)
}
