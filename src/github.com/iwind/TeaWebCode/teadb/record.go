package teadb

import (
	"reflect"
	"fmt"
	"regexp"
	"encoding/json"
	"strings"
)

type Record struct {
	Meta struct {
		Id         int64               `json:"id"`
		ModifiedAt int64               `json:"modifiedAt"`
		Indexes    map[string][]string `json:"indexes"`
	} `json:"meta"`
	Value map[string]interface{} `json:"value"`
}

func NewRecordFromJSON(jsonData []byte) (*Record, error) {
	record := &Record{}
	err := json.Unmarshal(jsonData, record)
	if err != nil {
		return nil, err
	}

	return record, nil
}

func (record *Record) Id() int64 {
	return record.Meta.Id
}

func (record *Record) FieldValue(field string) (interface{}, bool) {
	if record.Value == nil {
		return nil, false
	}

	fieldValue, found := record.Value[field]
	return fieldValue, found
}

func (record *Record) FieldEquals(value1 interface{}, value2 interface{}) bool {
	reflectValue := reflect.ValueOf(value2)
	valueType := reflectValue.Type()
	if valueType == nil {
		return value1 == nil
	}

	valueKind := valueType.Kind()
	if valueKind == reflect.Slice || valueKind == reflect.Array {
		countElements := reflectValue.Len()
		for i := 0; i < countElements; i ++ {
			value := reflectValue.Index(i).Interface()
			if value == value1 || strings.ToLower(record.formatString(value)) == strings.ToLower(record.formatString(value1)) {
				return true
			}
		}
		return false
	}

	return value1 == value2 || strings.ToLower(record.formatString(value2)) == strings.ToLower(record.formatString(value1))
}

func (record *Record) FieldContains(values interface{}, fieldValue interface{}) bool {
	reflectValue := reflect.ValueOf(values)
	countElements := reflectValue.Len()
	for i := 0; i < countElements; i ++ {
		value := reflectValue.Index(i).Interface()
		if value == fieldValue || fmt.Sprintf("%#v", value) == fmt.Sprintf("%#v", fieldValue) {
			return true
		}
	}
	return false
}

func (record *Record) MatchConds(conds []*QueryCond) bool {
	if len(conds) == 0 {
		return true
	}

	isValid := true
	for _, cond := range conds {
		fieldValue, found := record.FieldValue(cond.field)
		if !found {
			continue
		}
		if cond.operator == QueryOperatorAttr {
			if !record.FieldEquals(cond.value, fieldValue) {
				isValid = false
				continue
			}
		} else if cond.operator == QueryOperatorIn {
			if !record.FieldContains(cond.value, fieldValue) {
				isValid = false
				continue
			}
		} else if cond.operator == QueryOperatorMatch {
			if cond.value == nil {
				isValid = false
				continue
			}

			// @TODO 如果fieldValue是slice，则只要有一个元素匹配，就认为是匹配的
			if !cond.value.(*regexp.Regexp).MatchString(fmt.Sprintf("%#v", fieldValue)) {
				isValid = false
				continue
			}
		}
	}

	return isValid
}

func (record *Record) formatString(value interface{}) string {
	valueString, ok := value.(string)
	if ok {
		return valueString
	}
	return fmt.Sprintf("%#v", value)
}
