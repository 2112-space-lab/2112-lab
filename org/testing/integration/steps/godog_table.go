package steps

import (
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"strings"

	"github.com/cucumber/godog"
)

var (
	errGodogInvalidTable = errors.New("invalid scenario table")
)

func GodogTableToKeyValueMap[TKey, TVal ~string](table *godog.Table, tableHasHeader bool) (map[TKey]TVal, error) {
	res := make(map[TKey]TVal, len(table.Rows))
	for rNum, row := range table.Rows {
		if tableHasHeader && rNum == 0 {
			continue
		}
		if len(row.Cells) != 2 {
			return res, fmt.Errorf("key value table row has %d columns - expected %d [%w]", len(row.Cells), 2, errGodogInvalidTable)
		}
		key := row.Cells[0].Value
		value := row.Cells[1].Value
		if _, ok := res[TKey(key)]; ok {
			return res, fmt.Errorf("key value table has key [%s] defined more than once [%w]", key, errGodogInvalidTable)
		}
		res[TKey(key)] = TVal(value)
	}
	return res, nil
}

func GodogTableToSingleColumnSlice(table *godog.Table, tableHasHeader bool) ([]string, error) {
	res := make([]string, 0, len(table.Rows))
	for rNum, row := range table.Rows {
		if tableHasHeader && rNum == 0 {
			continue
		}
		if len(row.Cells) != 1 {
			return res, fmt.Errorf("key value table row has %d columns - expected %d [%w]", len(row.Cells), 2, errGodogInvalidTable)
		}
		value := row.Cells[0].Value
		res = append(res, value)
	}
	return res, nil
}

func GodogTableToSlice[T any](table *godog.Table) ([]T, error) {
	var tItem T
	res := make([]T, 0, len(table.Rows))
	dest := reflect.ValueOf(tItem)
	if dest.Kind() != reflect.Struct {
		return res, fmt.Errorf("invalid type parameter [%T] - expected a struct", tItem)
	}

	structPropFromHeader := map[int]string{}
	for i, row := range table.Rows {
		if i == 0 {
			for col, cell := range row.Cells {
				structPropFromHeader[col] = strings.ReplaceAll(cell.Value, " ", "")
			}
			continue
		}
		var item T
		itemReflect := reflect.ValueOf(&item).Elem()
		for j, cell := range row.Cells {
			if cell.Value == "" {
				continue
			}
			var propName = structPropFromHeader[j]
			f := itemReflect.FieldByName(propName)
			if !f.IsValid() {
				return res, fmt.Errorf("unknown field [%s] in type [%T]", propName, item)
			}
			if !f.CanSet() {
				return res, fmt.Errorf("unable to set field [%s] in type [%T]", propName, item)
			}

			if f.Kind() == reflect.Ptr {
				elemType := f.Type().Elem()
				ptrValue := reflect.New(elemType)
				base := ptrValue.Elem()

				switch elemType.Kind() {
				case reflect.String:
					base.SetString(cell.Value)
				case reflect.Bool:
					b, err := strconv.ParseBool(cell.Value)
					if err != nil {
						return res, fmt.Errorf("invalid bool value [%s] for field [%s] in type [%T]", cell.Value, propName, item)
					}
					base.SetBool(b)
				case reflect.Int, reflect.Int16, reflect.Int32, reflect.Int64:
					v, err := strconv.ParseInt(cell.Value, 10, 64)
					if err != nil {
						return res, fmt.Errorf("invalid int value [%s] for field [%s] in type [%T]", cell.Value, propName, item)
					}
					base.SetInt(v)
				case reflect.Float32, reflect.Float64:
					b, err := strconv.ParseFloat(cell.Value, 64)
					if err != nil {
						return res, fmt.Errorf("invalid float value [%s] for field [%s] in type [%T]", cell.Value, propName, item)
					}
					base.SetFloat(b)
				default:
					return res, fmt.Errorf("unsupported pointer type [%s] for field [%s] in type [%T]", elemType.Kind().String(), propName, item)
				}
				f.Set(ptrValue)
			} else {
				switch f.Kind() {
				case reflect.String:
					f.SetString(cell.Value)
				case reflect.Bool:
					b, err := strconv.ParseBool(cell.Value)
					if err != nil {
						return res, fmt.Errorf("invalid bool value [%s] for field [%s] in type [%T]", cell.Value, propName, item)
					}
					f.SetBool(b)
				case reflect.Int, reflect.Int16, reflect.Int32, reflect.Int64:
					v, err := strconv.ParseInt(cell.Value, 10, 64)
					if err != nil {
						return res, fmt.Errorf("invalid int value [%s] for field [%s] in type [%T]", cell.Value, propName, item)
					}
					f.SetInt(v)
				case reflect.Float32, reflect.Float64:
					b, err := strconv.ParseFloat(cell.Value, 64)
					if err != nil {
						return res, fmt.Errorf("invalid float value [%s] for field [%s] in type [%T]", cell.Value, propName, item)
					}
					f.SetFloat(b)
				default:
					return res, fmt.Errorf("unsupported type [%s] for field [%s] in type [%T]", f.Kind().String(), propName, item)
				}
			}
		}
		res = append(res, item)
	}
	return res, nil
}
