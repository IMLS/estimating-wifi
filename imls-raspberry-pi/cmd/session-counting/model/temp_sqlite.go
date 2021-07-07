package model

import (
	"errors"
	"fmt"
	"log"
	"os"
	"reflect"
	"strconv"
	"strings"

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
	db := TempDB{}
	// t := time.Now()
	// todaysDB := fmt.Sprintf("%v-%02d-%02d-%v.sqlite", t.Year(), int(t.Month()), int(t.Day()), name)
	// lw.Debug("db filename: %v", todaysDB)
	// filepath := filepath.Join(path, todaysDB)
	dbptr, err := sqlx.Open("sqlite3", path)
	if err != nil {
		lw.Debug("could not open temporary db: %v", path)
		lw.Fatal(err.Error())
	}

	db.DBName = name
	db.Path = path
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

	fields := make([]string, 0)
	for col, t := range columns {
		fp := fmt.Sprintf("%v %v", col, t)
		fields = append(fields, fp)
	}
	statement := fmt.Sprintf("CREATE TABLE IF NOT EXISTS %v (%v)", name, strings.Join(fields, ", "))
	lw.Debug("AddTable ", statement)
	_, err := tdb.Ptr.Exec(statement)
	if err != nil {
		lw.Info("could not re-create '", name, "' table in temporary db.")
		lw.Fatal(err.Error())
	}
}

func (tdb *TempDB) AddStructAsTable(table string, s interface{}) {
	//lw := logwrapper.NewLogger(nil)
	columns := make(map[string]string)
	rt := reflect.TypeOf(s)
	if rt.Kind() != reflect.Struct {
		log.Println("cannot add this struct as a table in", table, s)
		panic("bad type")
	}
	for i := 0; i < rt.NumField(); i++ {
		f := rt.Field(i)
		columns[f.Tag.Get("db")] = f.Tag.Get("sqlite")
	}

	tdb.AddTable(table, columns)
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
		return fmt.Sprintf("%v", v)
	case "date":
		// t, _ := time.Parse(time.RFC3339, fmt.Sprintf("%v", v))
		return fmt.Sprintf("%v", v)
	default:
		log.Fatal(fmt.Sprintf("could not convert type: %v, %v", t, v))

	}

	return v
}

func (tdb *TempDB) GetFields(table string) (fields []string) {
	for col, t := range tdb.Tables[table] {
		if !strings.Contains(t, "AUTOINCREMENT") {
			fields = append(fields, col)
		}
	}
	return fields
}

func (tdb *TempDB) InsertStruct(table string, s interface{}) {
	//lw := logwrapper.NewLogger(nil)
	values := make(map[string]interface{})
	rt := reflect.TypeOf(s)

	if rt.Kind() != reflect.Struct {
		log.Println("cannot add this struct as a table", table, s)
		panic("bad type")
	}
	// r := reflect.ValueOf(s)
	for i := 0; i < rt.NumField(); i++ {
		f := rt.Field(i)
		r := reflect.ValueOf(s)
		fv := reflect.Indirect(r).FieldByName(f.Name)
		//log.Println("v", fmt.Sprintf("%v", fv))
		tag := f.Tag
		// log.Println("tag", tag)
		// time.Sleep(1 * time.Second)
		if !strings.Contains(string(tag), "AUTOINCREMENT") {
			cleantag := strings.ReplaceAll(f.Tag.Get("db"), "\"", "")
			cleanvalue := strings.ReplaceAll(fmt.Sprintf("%v", fv), "\"", "")
			values[cleantag] = cleanvalue
		}
	}
	// log.Println("values", values)
	tdb.Insert(table, values)
}

func (tdb *TempDB) Insert(table string, values map[string]interface{}) {
	lw := logwrapper.NewLogger(nil)
	db := tdb.Ptr

	fields := make([]string, 0)
	subs := make([]interface{}, 0)
	questions := make([]string, 0)

	// log.Println("values", values)

	for col, v := range values {
		// Only process values that have matching columns in the table.
		if _, ok := tdb.Tables[table][col]; ok {
			t := tdb.Tables[table][col]
			fields = append(fields, col)
			subs = append(subs, convert(t, v))
			questions = append(questions, "?")
		}
	}

	// log.Println("subs", subs)

	full := fmt.Sprintf("INSERT INTO %v (%v) VALUES (%v)",
		table,
		strings.Join(fields, ", "),
		strings.Join(questions, ", "))

	// log.Println("full", full)

	insertS, err := db.Prepare(full)
	if err != nil {
		lw.Info("could not prepare %v insert statement", table)
		lw.Fatal(err.Error())
	}
	_, err = insertS.Exec(subs...)
	if err != nil {
		lw.Info("could not insert into temporary db: %v", table)
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
	r := make(map[string]interface{})
	for rows.Next() {
		rows.MapScan(r)
		log.Println(r)
	}
	return nil
}
