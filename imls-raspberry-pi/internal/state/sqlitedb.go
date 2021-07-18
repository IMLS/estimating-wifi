package state

import (
	"fmt"
	"log"
	"reflect"
	"strings"
	"sync"

	"github.com/jmoiron/sqlx"
	"gsa.gov/18f/internal/logwrapper"
)

type SqliteDB struct {
	Ptr  *sqlx.DB
	Path string
	// Tables map[string]map[string]string
	Tables map[string]*SqliteTable
	mutex  sync.Mutex
}

func NewSqliteDB(path string) *SqliteDB {
	// lw := logwrapper.NewLogger(nil)
	db := SqliteDB{}
	db.Path = path
	db.Ptr = nil
	//db.Tables = make(map[string]map[string]string)
	db.Tables = make(map[string]*SqliteTable)
	db.Open()
	return &db
}

func (db *SqliteDB) Open() {
	lw := logwrapper.NewLogger(nil)
	if db.Ptr == nil {
		// lw.Debug("opening db: ", tdb.DBName, " path: ", tdb.Path)
		dbptr, err := sqlx.Open("sqlite3", db.Path)
		if err != nil {
			lw.Error("could not open temporary db: ", db.Path)
			lw.Fatal(err.Error())
		} else {
			db.Ptr = dbptr
		}
	} else {
		lw.Debug("db already open: [ ", db.Path, " ]")
	}
}

func (db *SqliteDB) Close() {
	lw := logwrapper.NewLogger(nil)
	if strings.Contains(db.Path, "memory") {
		// Do nothing. Keep memory DB open.
	} else {
		if db.Ptr != nil {
			//lw.Debug("closing db: ", tdb.DBName)
			err := db.Ptr.Close()
			if err != nil {
				lw.Error("could not close db [", db.Path, "]")
			}
			db.Ptr = nil
		}
	}
}

func (db *SqliteDB) GetPtr() *sqlx.DB {
	return db.Ptr
}

type SqliteTable struct {
	Name            string
	DB              *sqlx.DB
	ColumnsAndTypes map[string]string
}

func (db *SqliteDB) CreateTable(name string) *SqliteTable {
	t := &SqliteTable{}
	t.Name = name
	t.ColumnsAndTypes = make(map[string]string)
	t.DB = db.Ptr
	return t
}

func (db *SqliteDB) CreateTableFromStruct(s interface{}) *SqliteTable {
	//columns := make(map[string]string)
	name := reflect.TypeOf(s).Name()
	t := db.CreateTable(name)
	ct := make(map[string]string)

	rt := reflect.TypeOf(s)
	if rt.Kind() != reflect.Struct {
		log.Println("cannot add this struct as a table in ", name, s)
		panic("bad type")
	}
	v := reflect.ValueOf(s)
	for i := 0; i < v.NumField(); i++ {
		f := reflect.TypeOf(s).Field(i)
		if f.Tag != "" {
			col := f.Tag.Get("db")
			tpe := f.Tag.Get("type")
			t.AddColumn(col, tpe)
			ct[col] = tpe
		}
	}
	cols := make([]string, 0)
	for c, t := range ct {
		cols = append(cols, fmt.Sprintf("%v %v", c, t))
	}
	stmnt := fmt.Sprintf("CREATE TABLE IF NOT EXISTS %v (%v)",
		t.Name+"s",
		strings.Join(cols, ", "))
	log.Println(stmnt)
	_, err := t.DB.Exec(stmnt)
	if err != nil {
		log.Fatalf("Failed to create table from struct: " + t.Name)
	}

	db.Tables[name] = t
	return t
}

func (db *SqliteDB) CheckTableExists(name string) bool {
	_, tableCheck := db.Ptr.Query("select * from " + name + ";")
	return tableCheck == nil
}

func (db *SqliteDB) GetTableFromStruct(s interface{}) *SqliteTable {
	name := reflect.TypeOf(s).Name()
	return db.Tables[name]
}

////////////////////////////////////////////////////////

func (t *SqliteTable) AddColumn(name string, sqlitetype string) {
	t.ColumnsAndTypes[name] = sqlitetype
}

func (t *SqliteTable) InsertStruct(s interface{}) {
	if reflect.ValueOf(s).Kind() == reflect.Struct {
		name := reflect.TypeOf(s).Name()
		columns := make([]string, 0)
		values := make([]string, 0)

		v := reflect.ValueOf(s)
		for i := 0; i < v.NumField(); i++ {
			f := reflect.TypeOf(s).Field(i)
			columns = append(columns, f.Name)
			switch v.Field(i).Kind() {
			case reflect.Int:
				values = append(values, fmt.Sprint(v.Field(i).Int()))
			case reflect.String:
				values = append(values, fmt.Sprintf("\"%v\"", v.Field(i).String()))
			default:
				fmt.Println("Unsupported field type in " + name)
				return
			}
		}
		query := fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s)",
			name+"s",
			strings.Join(columns, ", "),
			strings.Join(values, ", "))
		// fmt.Println(query)
		_, err := t.DB.Exec(query)
		if err != nil {
			log.Println("INSERT FAILED ON " + name)
			log.Println(query)
		}
	}

}

func (t *SqliteTable) GetIntegerType() string {
	return "INTEGER"
}

func (t *SqliteTable) GetTextType() string {
	return "TEXT"
}

func (t *SqliteTable) GetDateType() string {
	return "TEXT"
}
