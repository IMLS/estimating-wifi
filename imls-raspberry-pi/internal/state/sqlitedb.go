package state

import (
	"fmt"
	"log"
	"reflect"
	"strings"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
	"gsa.gov/18f/internal/interfaces"
	"gsa.gov/18f/internal/logwrapper"
)

type SqliteDB struct {
	Ptr    *sqlx.DB
	Path   string
	Tables map[string]*SqliteTable
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
	cfg := GetConfig()
	if db.Ptr == nil {
		// lw.Debug("opening db: ", tdb.DBName, " path: ", tdb.Path)
		ptr, err := sqlx.Open("sqlite3", db.Path)
		if err != nil {
			cfg.Log().Error("could not open temporary db: ", db.Path)
			cfg.Log().Fatal(err.Error())
		} else {
			db.Ptr = ptr
		}
	} else {
		// cfg.Log().Debug("db already open: [ ", db.Path, " ]")
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

func (db *SqliteDB) GetPath() string {
	return db.Path
}

func (db *SqliteDB) initTable(name string) *SqliteTable {
	t := &SqliteTable{}
	t.Name = name
	t.ColumnsAndTypes = make(map[string]string)
	t.DB = db
	return t
}

func (db *SqliteDB) InitTable(name string) interfaces.Table {
	t := db.initTable(name)
	db.Tables[name] = t
	return t
}

func (db *SqliteDB) RemoveTable(name string) {
	delete(db.Tables, name)
}

func (db *SqliteDB) CreateTableFromStruct(s interface{}) interfaces.Table {
	//columns := make(map[string]string)
	name := reflect.TypeOf(s).Name()
	t := db.initTable(name + "s")
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
			// log.Println("adding column " + f.Tag.Get("db"))
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
		t.Name,
		strings.Join(cols, ", "))
	// log.Println(stmnt)
	_, err := t.DB.GetPtr().Exec(stmnt)
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

func (db *SqliteDB) ListTables() []string {
	names := make([]string, 0)
	for name := range db.Tables {
		names = append(names, name)
	}
	return names
}

func (db *SqliteDB) GetTableFromStruct(s interface{}) interfaces.Table {
	name := reflect.TypeOf(s).Name()
	db.InitTable(name)
	// cfg.Log().Debug(db.Tables)
	return db.Tables[name]
}

func (db *SqliteDB) GetTableByName(name string) interfaces.Table {
	return db.Tables[name]
}

func (db *SqliteDB) Query(s string) (*sqlx.Rows, error) {
	return db.Ptr.Queryx(s)
}

////////////////////////////////////////////////////////

type SqliteTable struct {
	Name            string
	DB              interfaces.Database
	ColumnsAndTypes map[string]string
}

func (t *SqliteTable) AddColumn(name string, sqlitetype string) {
	// log.Println("adding column " + name + " type " + sqlitetype)
	t.ColumnsAndTypes[name] = sqlitetype
}

func (t *SqliteTable) Create() {
	cols := make([]string, 0)
	for c, t := range t.ColumnsAndTypes {
		cols = append(cols, fmt.Sprintf("%v %v", c, t))
	}
	stmnt := fmt.Sprintf("CREATE TABLE IF NOT EXISTS %v (%v)",
		t.Name,
		strings.Join(cols, ", "))
	// log.Println(stmnt)
	_, err := t.DB.GetPtr().Exec(stmnt)
	if err != nil {
		log.Println("Failed to create table: " + t.Name)
		log.Println("In DB " + t.DB.GetPath())
		log.Fatalf(err.Error())
	}
}

func (t *SqliteTable) InsertStruct(s interface{}) {
	if reflect.ValueOf(s).Kind() == reflect.Struct {
		name := reflect.TypeOf(s).Name()
		columns := make([]string, 0)
		values := make([]string, 0)

		v := reflect.ValueOf(s)
		for i := 0; i < v.NumField(); i++ {
			f := reflect.TypeOf(s).Field(i)
			if f.Tag != "" {
				col := f.Tag.Get("db")
				if !strings.Contains(f.Tag.Get("type"), "AUTOINCREMENT") {
					columns = append(columns, col)
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
			}
		}
		query := fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s)",
			name+"s",
			strings.Join(columns, ", "),
			strings.Join(values, ", "))
		// fmt.Println(query)
		_, err := t.DB.GetPtr().Exec(query)
		if err != nil {
			log.Println("INSERT FAILED ON " + name)
			log.Println(err.Error())
			log.Println(query)
		}
	}

}

func (t *SqliteTable) Drop() {
	ptr := t.DB.GetPtr()
	stmt := fmt.Sprintf("DROP TABLE IF EXISTS %v", t.Name)
	// log.Println(stmt)
	_, err := ptr.Exec(stmt)
	if err != nil {
		log.Fatal("could not drop table " + t.Name)
	}
	t.DB.RemoveTable(t.Name)
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

func (t *SqliteTable) GetDB() interfaces.Database {
	return t.DB
}
