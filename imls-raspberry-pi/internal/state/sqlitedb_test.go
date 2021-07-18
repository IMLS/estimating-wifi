package state

import (
	"log"
	"testing"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

type Apple struct {
	Color  string `db:"color" type:"TEXT"`
	Weight int    `db:"weight" type:"INTEGER"`
}

func (a Apple) SelectAll(db *sqlx.DB) []Apple {
	apples := []Apple{}
	err := db.Select(&apples, "SELECT * FROM Apples")
	if err != nil {
		log.Println("Error in SELECT *")
		log.Fatal(err.Error())
	}
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
	d := NewSqliteDB("/tmp/test.sqlite")
	t := d.CreateTable("apples")
	t.AddColumn("count", t.GetIntegerType())
}

func TestSqliteDB2(test *testing.T) {
	d := NewSqliteDB("/tmp/test.sqlite")
	t := d.CreateTableFromStruct(Apple{})
	t.InsertStruct(Apple{Color: "red", Weight: 3})
	t.InsertStruct(Apple{Color: "green", Weight: 5})
}
func TestSelectAll(test *testing.T) {
	d := NewSqliteDB("/tmp/test.sqlite")
	t := d.CreateTableFromStruct(Apple{})
	t.InsertStruct(Apple{Color: "red", Weight: 3})
	apples := Apple{}.SelectAll(t.DB)
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
}
