package teadb

import (
	"regexp"
	"github.com/syndtr/goleveldb/leveldb/util"
	"github.com/syndtr/goleveldb/leveldb/errors"
	"github.com/iwind/TeaGo/logs"
	"hash/crc32"
	"fmt"
	"bytes"
	"github.com/iwind/TeaGo/utils/number"
	"reflect"
	"strings"
)

const (
	QueryOperatorAttr  = "attr"
	QueryOperatorIn    = "in"
	QueryOperatorMatch = "match"
)

const (
	QuerySubActionFind   = "find"
	QuerySubActionDelete = "delete"
	QuerySubActionCount  = "count"
)

type Query struct {
	conds    map[string]map[string][]*QueryCond // { field => { operator => [ cond1, cond2, .. ]  }  }
	offset   int64
	limit    int64
	db       *DB
	dataType string
	reverse  bool

	subAction string
}

type QueryCond struct {
	field    string
	operator string
	value    interface{}
	append   bool
}

func NewQuery(db *DB, dataType string) *Query {
	return &Query{
		conds:    map[string]map[string][]*QueryCond{},
		db:       db,
		dataType: dataType,
		offset:   0,
		limit:    -1,
	}
}

func (query *Query) Reverse(bool bool) *Query {
	query.reverse = bool
	return query
}

func (query *Query) Cond(cond *QueryCond) *Query {
	fieldMap, found := query.conds[cond.field]
	if !found {
		fieldMap = map[string][]*QueryCond{}
	}

	conds, found := fieldMap[cond.operator]
	if !found {
		conds = []*QueryCond{}
	}
	if cond.append {
		conds = append(conds, cond)
	} else {
		conds = []*QueryCond{cond}
	}
	fieldMap[cond.operator] = conds
	query.conds[cond.field] = fieldMap
	return query
}

func (query *Query) Attr(field string, value interface{}) *Query {
	cond := &QueryCond{
		field:    field,
		operator: QueryOperatorAttr,
		value:    value,
		append:   false,
	}
	query.Cond(cond)
	return query
}

func (query *Query) In(field string, values interface{}) *Query {
	valuesKind := reflect.TypeOf(values).Kind()
	if valuesKind != reflect.Slice && valuesKind != reflect.Array {
		logs.Errorf("Query.In(field, values): values must be slice or array")
		return query
	}

	cond := &QueryCond{
		field:    field,
		operator: QueryOperatorIn,
		value:    values,
		append:   false,
	}
	query.Cond(cond)
	return query
}

func (query *Query) Match(field string, regexp *regexp.Regexp) *Query {
	cond := &QueryCond{
		field:    field,
		operator: QueryOperatorMatch,
		value:    regexp,
		append:   true,
	}
	query.Cond(cond)
	return query
}

func (query *Query) Offset(offset int64) *Query {
	query.offset = offset
	return query
}

func (query *Query) Limit(size int64) *Query {
	query.limit = size
	return query
}

func (query *Query) run(result interface{}) error {
	var count = int64(0)

	offset := query.offset
	if offset < 0 {
		offset = 0
	}

	if query.limit == 0 {
		return nil
	}

	var flatConds = []*QueryCond{}
	var firstCond *QueryCond
	if len(query.conds) > 0 {
		for _, operators := range query.conds {
			for _, conds := range operators {
				for _, cond := range conds {
					//@TODO 选择最少匹配记录的条件
					// @TODO 支持 QueryOperatorIn
					if firstCond == nil && cond.operator == QueryOperatorAttr {
						firstCond = cond
					} else {
						flatConds = append(flatConds, cond)
					}
				}
			}
		}
	}

	var iteratorStarted = false
	if firstCond == nil {
		prefix := query.dataType + ".DATA."
		it := query.db.native.NewIterator(util.BytesPrefix([]byte(prefix)), nil)

		for {
			if query.reverse {
				if iteratorStarted {
					if !it.Prev() {
						break
					}
				} else {
					if !it.Last() {
						break
					}
					iteratorStarted = true
				}
			} else {
				if !it.Next() {
					break
				}
			}

			if len(flatConds) == 0 && query.subAction == QuerySubActionCount {
				count ++

				if count <= offset {
					continue
				}

				*result.(*int64) ++

				if query.limit > 0 && count >= offset+query.limit {
					break
				}
			} else {
				record, err := NewRecordFromJSON(it.Value())
				if err != nil {
					logs.Error(err)
					continue
				}

				if record.MatchConds(flatConds) {
					count ++

					if count <= offset {
						continue
					}

					if query.subAction == QuerySubActionFind {
						*result.(*[]*Record) = append(*result.(*[]*Record), record)
					} else if query.subAction == QuerySubActionDelete {
						err := query.db.Delete(query.dataType, record.Id())
						if err != nil {
							return err
						}
						*result.(*int64) ++
					} else if query.subAction == QuerySubActionCount {
						*result.(*int64) ++
					}

					if query.limit > 0 && count >= offset+query.limit {
						break
					}
				}
			}
		}
	} else if firstCond.operator == QueryOperatorAttr {
		valueString, ok := firstCond.value.(string)
		if !ok {
			valueString = fmt.Sprintf("%#v", firstCond.value)
		}
		valueString = strings.ToLower(valueString)

		prefix := query.dataType + ".INDEX.[" + firstCond.field + "]." + fmt.Sprintf("%d", crc32.ChecksumIEEE([]byte(valueString))) + "."
		it := query.db.native.NewIterator(util.BytesPrefix([]byte(prefix)), nil)

		for {
			if query.reverse {
				if iteratorStarted {
					if !it.Prev() {
						break
					}
				} else {
					if !it.Last() {
						break
					}
					iteratorStarted = true
				}
			} else {
				if !it.Next() {
					break
				}
			}

			if len(flatConds) == 0 && query.subAction == QuerySubActionCount {
				count ++

				if count <= offset {
					continue
				}

				*result.(*int64) ++

				if query.limit > 0 && count >= offset+query.limit {
					break
				}
			} else {
				id := numberutil.ParseInt64(string(it.Key()[bytes.LastIndex(it.Key(), []byte("."))+1:]), 0)
				record, err := query.db.Get(query.dataType, id)
				if err != nil {
					return err
				}

				if record == nil {
					continue
				}

				if record.MatchConds(flatConds) {
					count ++

					if count <= offset {
						continue
					}

					if query.subAction == QuerySubActionFind {
						*result.(*[]*Record) = append(*result.(*[]*Record), record)
					} else if query.subAction == QuerySubActionDelete {
						err := query.db.Delete(query.dataType, record.Id())
						if err != nil {
							return err
						}
						*result.(*int64) ++
					}

					if query.limit > 0 && count >= offset+query.limit {
						break
					}
				}
			}
		}
	}

	return nil
}

func (query *Query) FindAll() ([]*Record, error) {
	result := []*Record{}
	query.subAction = QuerySubActionFind
	err := query.run(&result)
	return result, err
}

func (query *Query) Find() (*Record, error) {
	ones, err := query.Offset(0).Limit(1).FindAll()
	if err != nil {
		return nil, err
	}

	if len(ones) == 0 {
		return nil, errors.New("not found")
	}

	return ones[0], nil
}

func (query *Query) FindField(field string) (interface{}, error) {
	one, err := query.Find()
	if err != nil {
		return nil, err
	}

	value, _ := one.FieldValue(field)
	return value, nil
}

func (query *Query) Delete() (int64, error) {
	query.subAction = QuerySubActionDelete
	count := int64(0)
	err := query.run(&count)
	return count, err
}

func (query *Query) Count() (int64, error) {
	query.subAction = QuerySubActionCount
	count := int64(0)
	err := query.run(&count)
	return count, err
}
