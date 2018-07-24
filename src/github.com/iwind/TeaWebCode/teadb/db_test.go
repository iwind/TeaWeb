package teadb

import (
	"testing"
	"github.com/iwind/TeaGo/Tea"
	"github.com/iwind/TeaGo/utils/time"
	"time"
	"github.com/iwind/TeaGo/logs"
)

func TestDB_NextId(t *testing.T) {
	db, err := NewDB(Tea.TmpDir() + Tea.Ds + "teadb")
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	nextId, err := db.nextID()
	if err != nil {
		t.Fatal(err)
	}
	t.Log(nextId)
}

func TestDB_Put(t *testing.T) {
	db, err := NewDB(Tea.TmpDir() + Tea.Ds + "teadb")
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	names := []string{"liu", "chao", "lu", "ping", "joe"}
	ages := []int{20, 18, 19, 20, 22}
	books := [][]string{{"golang", "php"}, {"golang", "python", "php"}, {"php"}, {"python", "golang"}, {"golang"}}
	for index, name := range names {
		_, err := db.Put("LOG", map[string]interface{}{
			"name":  name,
			"age":   ages[index],
			"time":  timeutil.Format(time.Now(), "H:i:s"),
			"books": books[index],
		})

		if err != nil {
			t.Fatal(err)
		}
	}
}

func TestDB_PutStruct(t *testing.T) {
	db, err := NewDB(Tea.TmpDir() + Tea.Ds + "teadb")
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	s := &struct {
		Name  string   `json:"name"`
		Age   int      `json:"age"`
		Books []string `json:"books"`
	}{
		Name:  "liu",
		Age:   28,
		Books: []string{"PHP"},
	}

	t.Log(db.PutStruct("LOG", s))
}

func TestDB_Get(t *testing.T) {
	db, err := NewDB(Tea.TmpDir() + Tea.Ds + "teadb")
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	t.Log(db.Get("LOG", 11))

	one, err := db.Get("LOG", 24)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(one.Value, one.Id())

	value, err := db.native.Get([]byte("LOG.INDEX.[age].2944839123.00000000000000000000000000000024"), nil)
	if err != nil {
		t.Fatal(err)
	}

	t.Log(string(value))
}

func TestDb_createIndex(t *testing.T) {
	db, err := NewDB(Tea.TmpDir() + Tea.Ds + "teadb")
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	indexes := map[string]*IndexInfo{}
	value := map[string]interface{}{
		"name":      "liu",
		"age":       20,
		"gender":    "boy",
		"isWorking": true,
		"house": map[string]interface{}{
			"price": 10000,
			"size":  105.5,
			"location": map[string]interface{}{
				"lat":     39.0123,
				"lng":     106.447,
				"address": "Fun Hill",
			},
		},
		"books":  []string{"Golang", "PHP", "Python"},
		"orders": []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10},
	}
	for key, fieldValue := range value {
		db.createIndexes(indexes, "LOG", 123, fieldValue, key)
	}

	logs.Dump(indexes)
}

func TestDB_Delete(t *testing.T) {
	db, err := NewDB(Tea.TmpDir() + Tea.Ds + "teadb")
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	t.Log(db.Delete("LOG", 24))
}
