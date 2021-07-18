package interfaces

import (
	"github.com/jmoiron/sqlx"
)

type Database interface {
	// NewDB(path string) *sqlx.DB
	Open()
	Close()
	GetPtr() *sqlx.DB
	CreateTable(name string) Table
	CreateTableFromStruct(s interface{}) Table
	CheckTableExists(name string) bool
	ListTables() []string
	//GetTable(name string) Table
	GetTableFromStruct(s interface{}) Table
	Exec() error
	Query()
}

type Table interface {
	AddColumn(name string, sqltype string)
	// AddColumns(namestypes map[string]string)
	InsertStruct(s interface{})

	GetIntegerType() string
	GetTextType() string
	GetDateType() string
	//FieldToDate(field string) time.Time
	//DateToField(t time.Time) string
}
