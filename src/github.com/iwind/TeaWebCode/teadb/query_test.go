package teadb

import (
	"testing"
	"github.com/iwind/TeaGo/Tea"
	"regexp"
	"time"
)

func TestNewQuery(t *testing.T) {
	db, err := NewDB(Tea.TmpDir() + Tea.Ds + "teadb")
	if err != nil {
		t.Fatal(err)
	}

	query := db.NewQuery("LOG")
	t.Log(query.Reverse(false).Offset(3).Limit(5).FindAll())
}

func TestNewQueryReverse(t *testing.T) {
	db, err := NewDB(Tea.TmpDir() + Tea.Ds + "teadb")
	if err != nil {
		t.Fatal(err)
	}

	query := db.NewQuery("LOG")
	ones, err := query.Reverse(true).Offset(0).Limit(5).FindAll()
	if err != nil {
		t.Fatal(err)
	}
	for _, one := range ones {
		t.Log(one.Id(), one.Value)
	}
}

func TestQuery_Attr(t *testing.T) {
	db, err := NewDB(Tea.TmpDir() + Tea.Ds + "teadb")
	if err != nil {
		t.Fatal(err)
	}

	timeBefore := time.Now()
	defer func() {
		t.Log(float64(time.Since(timeBefore).Nanoseconds())/1000000, "ms")

		defer db.Close()
	}()

	query := db.NewQuery("LOG")

	reg, err := regexp.Compile("(?i)o")
	if err != nil {
		t.Fatal(err)
	}

	ones, err := query.Reverse(false).
		Attr("books", "golang").
	//Attr("name", "liu").
	//In("age", [...]int{18, 19, 20, 21, 22}).

		Match("name", reg).

		Offset(0).
		Limit(2).
		Reverse(true).
		FindAll()
	if err != nil {
		t.Fatal(err)
	}
	for _, one := range ones {
		t.Log(one.Id(), one.Value)
	}
}

func TestQuery_AttrCase(t *testing.T) {
	db, err := NewDB(Tea.TmpDir() + Tea.Ds + "teadb")
	if err != nil {
		t.Fatal(err)
	}

	timeBefore := time.Now()
	defer func() {
		t.Log(float64(time.Since(timeBefore).Nanoseconds())/1000000, "ms")

		defer db.Close()
	}()

	query := db.NewQuery("LOG")

	ones, err := query.Reverse(false).
		Attr("books", "PHP").
		Attr("name", "lu").
		Attr("age", 19).

		Offset(0).
		Limit(10).
		Reverse(true).
		FindAll()
	if err != nil {
		t.Fatal(err)
	}
	for _, one := range ones {
		t.Log(one.Id(), one.Value)
	}
}

func TestQuery_Find(t *testing.T) {
	db, err := NewDB(Tea.TmpDir() + Tea.Ds + "teadb")
	if err != nil {
		t.Fatal(err)
	}

	timeBefore := time.Now()
	defer func() {
		t.Log(float64(time.Since(timeBefore).Nanoseconds())/1000000, "ms")

		defer db.Close()
	}()

	record, err := db.NewQuery("LOG").Attr("books", "golang").Find()
	if err != nil {
		t.Fatal(err)
	} else {
		t.Log(record.Value)
	}
}

func TestQuery_FindField(t *testing.T) {
	db, err := NewDB(Tea.TmpDir() + Tea.Ds + "teadb")
	if err != nil {
		t.Fatal(err)
	}

	timeBefore := time.Now()
	defer func() {
		t.Log(float64(time.Since(timeBefore).Nanoseconds())/1000000, "ms")

		defer db.Close()
	}()

	value, err := db.NewQuery("LOG").Attr("books", "golang").FindField("name")
	if err != nil {
		t.Fatal(err)
	} else {
		t.Log(value)
	}
}

func TestQuery_Delete(t *testing.T) {
	db, err := NewDB(Tea.TmpDir() + Tea.Ds + "teadb")
	if err != nil {
		t.Fatal(err)
	}

	timeBefore := time.Now()
	defer func() {
		t.Log(float64(time.Since(timeBefore).Nanoseconds())/1000000, "ms")

		defer db.Close()
	}()

	t.Log(db.NewQuery("LOG").Attr("books", "golang").Delete())
}

func TestQuery_Count(t *testing.T) {
	db, err := NewDB(Tea.TmpDir() + Tea.Ds + "teadb")
	if err != nil {
		t.Fatal(err)
	}

	timeBefore := time.Now()
	defer func() {
		t.Log(float64(time.Since(timeBefore).Nanoseconds())/1000000, "ms")

		defer db.Close()
	}()

	t.Log(db.NewQuery("LOG").Attr("books", "golang").Offset(0).Limit(10).Count())
}
