package model

import (
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
	"gsa.gov/18f/logwrapper"
)

type TempDB struct {
	Ptr    *sqlx.DB
	DBName string
	Path   string
	Tables map[string]map[string]string
}

// func newTempDbInFS(cfg *config.Config) *sqlx.DB {
// 	lw := logwrapper.NewLogger(nil)

// 	t := time.Now()
// 	todaysDB := fmt.Sprintf("%v-%02d-%02d-wifi.sqlite", t.Year(), int(t.Month()), int(t.Day()))
// 	lw.Info("Created temporary db: %v", todaysDB)
// 	path := filepath.Join(cfg.Local.WebDirectory, todaysDB)
// 	db, err := sqlx.Open("sqlite3", path)
// 	if err != nil {
// 		lw.Fatal("could not open temporary db: %v", path)
// 	}

// 	createWifiTable(cfg, db)
// 	return db
// }

func NewSqliteDB(name string, path string) *TempDB {
	lw := logwrapper.NewLogger(nil)
	if lw == nil {
		log.Println("lw is nil here...")
	}
	db := TempDB{}
	t := time.Now()
	todaysDB := fmt.Sprintf("%v-%02d-%02d-%v.sqlite", t.Year(), int(t.Month()), int(t.Day()), name)
	lw.Debug("db filename: %v", todaysDB)
	filepath := filepath.Join(path, todaysDB)
	dbptr, err := sqlx.Open("sqlite3", filepath)
	if err != nil {
		lw.Debug("could not open temporary db: %v", filepath)
		lw.Fatal(err.Error())
	}

	db.DBName = name
	db.Path = filepath
	db.Ptr = dbptr
	db.Tables = make(map[string]map[string]string)
	return &db
}

func (tdb *TempDB) DropTable(name string) {
	lw := logwrapper.NewLogger(nil)
	if _, ok := tdb.Tables[name]; ok {
		delete(tdb.Tables, name)
		stmt := fmt.Sprintf("DROP TABLE %v", name)
		_, err := tdb.Ptr.Exec(stmt)
		if err != nil {
			lw.Error("Could not drop table %v", name)
		}
	}
}

func (tdb *TempDB) AddTable(name string, columns map[string]string) {
	lw := logwrapper.NewLogger(nil)
	tdb.Tables[name] = columns

	columnstring := ""
	for col, t := range columns {
		if columnstring == "" {
			columnstring = fmt.Sprintf("\"%v\" %v", col, t)
		} else {
			columnstring = fmt.Sprintf("%v, \"%v\" %v", columnstring, col, t)
		}
	}

	statement := fmt.Sprintf("CREATE TABLE IF NOT EXISTS %v (%v)", name, columnstring)

	_, err := tdb.Ptr.Exec(statement)
	if err != nil {
		lw.Info("could not re-create %v table in temporary db.", name)
		lw.Fatal(err.Error())
	}
}

func convert(t string, v interface{}) interface{} {
	toggle := ""
	searches := []string{"integer", "text", "date"}

	for _, s := range searches {
		if strings.Contains(strings.ToLower(t), s) {
			toggle = s
		}
	}

	switch toggle {
	case "integer":
		i, _ := strconv.Atoi(fmt.Sprintf("%v", v))
		return i
	case "text":
		return v
	case "date":
		t, _ := time.Parse(time.RFC3339, fmt.Sprintf("%v", v))
		return t
	default:
		log.Fatal(fmt.Sprintf("could not convert type: %v, %v", t, v))

	}

	return v
}

func (tdb *TempDB) GetFields(table string) (fields []string) {
	for col, t := range tdb.Tables[table] {
		if !strings.Contains(t, "PRIMARY") {
			fields = append(fields, col)
		}
	}
	return fields
}

func (tdb *TempDB) Insert(name string, values map[string]interface{}) {
	lw := logwrapper.NewLogger(nil)
	db := tdb.Ptr

	insertstatement := fmt.Sprintf("INSERT INTO %v ", name)
	valuestoinsert := make([]interface{}, 0)
	numcols := 0
	questions := "(?"
	for col, v := range values {
		// Only process values that have matching columns in the table.
		if _, ok := tdb.Tables[name][col]; ok {
			t := tdb.Tables[name][col]
			if numcols == 0 {
				insertstatement = fmt.Sprintf("%v(%v", insertstatement, col)
				valuestoinsert = append(valuestoinsert, convert(t, v))
			} else if numcols == len(values)-1 {
				insertstatement = fmt.Sprintf("%v, %v) ", insertstatement, col)
				questions = fmt.Sprintf("%v, ?)", questions)
				valuestoinsert = append(valuestoinsert, convert(t, v))
			} else {
				insertstatement = fmt.Sprintf("%v, %v", insertstatement, col)
				questions = fmt.Sprintf("%v, ?", questions)
				valuestoinsert = append(valuestoinsert, convert(t, v))
			}
		}
		numcols += 1
	}

	full := fmt.Sprintf("%v VALUES %v", insertstatement, questions)
	insertS, err := db.Prepare(full)
	if err != nil {
		lw.Info("could not prepare %v insert statement", name)
		lw.Fatal(err.Error())
	}
	_, err = insertS.Exec(valuestoinsert...)
	if err != nil {
		lw.Info("could not insert into temporary db: %v", name)
		lw.Fatal(err.Error())
	}
}

func (tdb *TempDB) Close() {
	tdb.Ptr.Close()
}

func (tdb *TempDB) Remove() {
	lw := logwrapper.NewLogger(nil)
	err := os.Remove(tdb.Path)
	if err != nil {
		lw.Error("could not delete file: %v", tdb.Path)
	}
}

func (tdb *TempDB) DebugDump(name string) error {
	lw := logwrapper.NewLogger(nil)
	q := fmt.Sprintf("SELECT * FROM %v", name)
	rows, err := tdb.Ptr.Queryx(q)
	if err != nil {
		lw.Info("could not select all from %v", name)
		return errors.New("could not select all from db")
	}
	for rows.Next() {
		var r map[string]string
		rows.Scan(r)
		log.Println(r)
	}
	return nil
}

func (tdb *TempDB) SelectAll(name string, arr interface{}) {
	lw := logwrapper.NewLogger(nil)
	err := tdb.Ptr.Select(&arr, fmt.Sprintf("SELECT * FROM %v", name))
	if err != nil {
		lw.Info("error in extracting all events: %v", name)
		lw.Fatal(err.Error())
	}

	lw.Length(name, arr)
}
