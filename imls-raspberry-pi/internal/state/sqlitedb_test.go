package state

import (
	"log"
	"os"
	"testing"
	"time"

	"net/http"
	_ "net/http/pprof"

	_ "github.com/mattn/go-sqlite3"
	"gsa.gov/18f/internal/interfaces"
)

type Apple struct {
	Color  string `db:"color" type:"TEXT"`
	Weight int    `db:"weight" type:"INTEGER"`
}

func (a Apple) SelectAll(db interfaces.Database) []Apple {
	apples := []Apple{}
	db.GetPtr().Select(&apples, "SELECT * FROM Apples")
	// if err != nil {
	// 	log.Println("Found no apples")
	// 	log.Println(err.Error())
	// }
	return apples
}

func AsApples(is []interface{}) []Apple {
	apples := make([]Apple, 0)
	for _, i := range is {
		apples = append(apples, i.(Apple))
	}
	return apples
}

func TestSqliteDB(test *testing.T) {
	tempDB, err := os.CreateTemp("", "sqlitedb-test-db")
	if err != nil {
		log.Fatal(err)
	}
	defer os.Remove(tempDB.Name())
	d := NewSqliteDB(tempDB.Name())
	t := d.InitTable("oranges")
	t.AddColumn("count", t.GetIntegerType())
	t.Create()
}

func TestSqliteDB2(test *testing.T) {
	tempDB, err := os.CreateTemp("", "sqlitedb-test-db2")
	if err != nil {
		log.Fatal(err)
	}
	defer os.Remove(tempDB.Name())
	d := NewSqliteDB(tempDB.Name())
	t := d.CreateTableFromStruct(Apple{})
	t.InsertStruct(Apple{Color: "red", Weight: 3})
	t.InsertStruct(Apple{Color: "green", Weight: 5})
}

func TestSelectAll(test *testing.T) {
	tempDB, err := os.CreateTemp("", "sqlitedb-test-select-all")
	if err != nil {
		log.Fatal(err)
	}
	defer os.Remove(tempDB.Name())
	d := NewSqliteDB(tempDB.Name())
	t := d.CreateTableFromStruct(Apple{})
	t.InsertStruct(Apple{Color: "red", Weight: 3})
	apples := Apple{}.SelectAll(t.GetDB())
	if len(apples) < 0 {
		log.Fatal("no apples found")
	}
	isred := false
	for _, a := range apples {
		if a.Color == "red" {
			isred = true
		}
	}
	if !isred {
		log.Fatal("no red apples found")
	}
	t.Drop()
	apples = Apple{}.SelectAll(t.GetDB())

	if len(apples) != 0 {
		log.Println(apples)
		log.Fatal("found apples on a dropped table?")
	}
}

func TestManyOpens(test *testing.T) {
	tempDB, err := os.CreateTemp("", "sqlitedb-test-many-opens")
	if err != nil {
		log.Fatal(err)
	}
	defer os.Remove(tempDB.Name())
	go http.ListenAndServe("localhost:8080", nil)
	for i := 0; i < 8000; i++ {
		d := NewSqliteDB(tempDB.Name())
		d.CreateTableFromStruct(Apple{})
	}
	time.Sleep(10 * time.Second)
	for i := 0; i < 16000; i++ {
		d := NewSqliteDB(tempDB.Name())
		d.CreateTableFromStruct(Apple{})
	}
}
