package state

import (
	"log"
	"testing"

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
	d := NewSqliteDB("/tmp/test.sqlite")
	t := d.InitTable("oranges")
	t.AddColumn("count", t.GetIntegerType())
	t.Create()
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
