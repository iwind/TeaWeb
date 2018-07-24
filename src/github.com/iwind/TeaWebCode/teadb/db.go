package teadb

import (
	"github.com/syndtr/goleveldb/leveldb"
	"github.com/iwind/TeaGo/logs"
	"github.com/pquerna/ffjson/ffjson"
	"github.com/syndtr/goleveldb/leveldb/errors"
	"sync"
	"fmt"
	"time"
	"hash/crc32"
	"reflect"
	"strings"
)

type DB struct {
	native *leveldb.DB
	locker *sync.Mutex
}

type GlobalMeta struct {
	Id int64 `json:"id"`
}

type IndexInfo struct {
	key    string
	subKey string
	value  string
}

func NewDB(path string) (*DB, error) {
	db := &DB{
		locker: &sync.Mutex{},
	}
	err := db.Open(path)
	return db, err
}

func (db *DB) Open(path string) error {
	native, err := leveldb.OpenFile(path, nil)
	if err != nil {
		logs.Errorf("db.Open(path): %s", err.Error())
		return err
	}

	db.native = native
	return nil
}

func (db *DB) Close() error {
	if db.native != nil {
		return db.native.Close()
	}
	return nil
}

func (db *DB) Put(dataType string, value map[string]interface{}) (id int64, err error) {
	nextId, err := db.nextID()
	if err != nil {
		return nextId, err
	}

	indexes := map[string]*IndexInfo{}
	indexMap := map[string][]string{} // key => [ subKey1, subKey2, ... ]
	for itemKey, itemValue := range value {
		db.createIndexes(indexes, dataType, nextId, itemValue, itemKey)
	}

	for indexKey, index := range indexes {
		arr, found := indexMap[index.key]
		if !found {
			arr = []string{index.subKey}
		} else {
			arr = append(arr, index.subKey)
		}
		indexMap[index.key] = arr

		err := db.native.Put([]byte(indexKey), []byte(index.value), nil)
		if err != nil {
			// @TODO 删除已存储的索引

			return nextId, err
		}
	}

	record := &Record{
		Meta: struct {
			Id         int64               `json:"id"`
			ModifiedAt int64               `json:"modifiedAt"`
			Indexes    map[string][]string `json:"indexes"`
		}{
			Id:         nextId,
			ModifiedAt: time.Now().UnixNano(),
			Indexes:    indexMap,
		},
		Value: value,
	}
	jsonData, err := ffjson.Marshal(record)
	if err != nil {
		return nextId, err
	}

	newKey := []byte(dataType + ".DATA." + db.formatId(nextId))
	return nextId, db.native.Put(newKey, jsonData, nil)
}

func (db *DB) PutStruct(dataType string, value interface{}) (id int64, err error) {
	data, err := ffjson.Marshal(value)
	if err != nil {
		return 0, err
	}

	dataMap := map[string]interface{}{}
	err = ffjson.Unmarshal(data, &dataMap)
	if err != nil {
		return 0, err
	}

	return db.Put(dataType, dataMap)
}

func (db *DB) Delete(dataType string, id int64) error {
	record, err := db.Get(dataType, id)
	if err != nil {
		if err == leveldb.ErrNotFound {
			return nil
		}
		return err
	}

	idString := db.formatId(id)
	err = db.native.Delete([]byte(dataType+".DATA."+idString), nil)
	if err != nil {
		return err
	}

	for key, subKeys := range record.Meta.Indexes {
		for _, subKey := range subKeys {
			fullKey := dataType + ".INDEX.[" + key + "]." + subKey
			err := db.native.Delete([]byte(fullKey), nil)
			if err != nil {
				logs.Error(err)
				return err
			}
		}
	}

	return nil
}

func (db *DB) createIndexes(indexes map[string]*IndexInfo, dataType string, id int64, value interface{}, key string) {
	isFormatted := false
	valueString := ""
	idString := db.formatId(id)

	valueType := reflect.TypeOf(value)
	valueKind := "nil"
	if valueType == nil {
		valueString = ""
		isFormatted = true
	} else {
		valueKind = valueType.Kind().String()
		switch valueKind {
		case "string":
			isFormatted = true
			valueString = value.(string)
		case "int":
			isFormatted = true
			valueString = fmt.Sprintf("%d", value)
		case "int8":
			isFormatted = true
			valueString = fmt.Sprintf("%d", value)
		case "int16":
			isFormatted = true
			valueString = fmt.Sprintf("%d", value)
		case "int32":
			isFormatted = true
			valueString = fmt.Sprintf("%d", value)
		case "int64":
			isFormatted = true
			valueString = fmt.Sprintf("%d", value)
		case "uint":
			isFormatted = true
			valueString = fmt.Sprintf("%d", value)
		case "uint8":
			isFormatted = true
			valueString = fmt.Sprintf("%d", value)
		case "uint16":
			isFormatted = true
			valueString = fmt.Sprintf("%d", value)
		case "uint32":
			isFormatted = true
			valueString = fmt.Sprintf("%d", value)
		case "uint64":
			isFormatted = true
			valueString = fmt.Sprintf("%d", value)
		case "bool":
			isFormatted = true
			if value.(bool) {
				valueString = "true"
			} else {
				valueString = "false"
			}
		case "float32":
			isFormatted = true
			valueString = fmt.Sprintf("%f", value)
		case "float64":
			isFormatted = true
			valueString = fmt.Sprintf("%f", value)
		case "map":
			valueMap, ok := value.(map[string]interface{})
			if ok {
				for itemKey, itemValue := range valueMap {
					db.createIndexes(indexes, dataType, id, itemValue, key+"."+itemKey)
				}
			}
			return
		case "slice":
			reflectValue := reflect.ValueOf(value)
			for i := 0; i < reflectValue.Len(); i ++ {
				db.createIndexes(indexes, dataType, id, reflectValue.Index(i).Interface(), key)
			}
		case "array":
			reflectValue := reflect.ValueOf(value)
			for i := 0; i < reflectValue.Len(); i ++ {
				db.createIndexes(indexes, dataType, id, reflectValue.Index(i).Interface(), key)
			}
		default:
			return
		}
	}

	if isFormatted {
		valueCRC := fmt.Sprintf("%d", crc32.ChecksumIEEE([]byte(strings.ToLower(valueString))))
		subIndexName := valueCRC + "." + idString
		indexName := dataType + ".INDEX.[" + key + "]." + subIndexName

		indexes[indexName] = &IndexInfo{
			key:    key,
			subKey: subIndexName,
			value:  valueKind + ":" + fmt.Sprintf("%d", id) + ":" + valueString,
		}
		return
	}
}

func (db *DB) Get(dataType string, id int64) (*Record, error) {
	newKey := []byte(dataType + ".DATA." + db.formatId(id))
	value, err := db.native.Get(newKey, nil)
	if err != nil {
		return nil, err
	}

	return NewRecordFromJSON(value)
}

func (db *DB) NewQuery(dataType string) *Query {
	return NewQuery(db, dataType)
}

func (db *DB) nextID() (int64, error) {
	// @TODO 将全局的meta分解为多个记录
	value, err := db.native.Get([]byte("$"), nil)
	if err != nil {
		if err != errors.ErrNotFound {
			return 0, err
		} else {
			value = []byte(`{ "id": 0 }`)
		}
	}

	db.locker.Lock()
	defer db.locker.Unlock()

	meta := &GlobalMeta{}
	err = ffjson.Unmarshal(value, meta)
	if err != nil {
		return 0, err
	}

	meta.Id += 1
	metaData, err := ffjson.Marshal(meta)
	if err != nil {
		return 0, err
	}

	err = db.native.Put([]byte("$"), metaData, nil)
	if err != nil {
		return 0, err
	}

	return meta.Id, nil
}

func (db *DB) formatId(id int64) string {
	return fmt.Sprintf("%032d", id)
}
