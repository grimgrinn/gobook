package params

import (
	"fmt"
	"net/url"
	"reflect"
	"strconv"
	"strings"
)

func Pack(ptr interface{}) (string, error) {
	v := reflect.ValueOf(ptr).Elem()
	t := v.Type()

	params := url.Values{}

	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		fieldType := t.Field(i)

		name := fieldType.Tag.Get("http")
		if name == "" {
			name = strings.ToLower(fieldType.Name)
		}

		if field.IsZero() {
			continue
		}

		switch field.Kind() {
		case reflect.String:
			params.Add(name, field.String())

		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			params.Add(name, strconv.FormatInt(field.Int(), 10))

		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			params.Add(name, strconv.FormatUint(field.Uint(), 10))

		case reflect.Float32, reflect.Float64:
			params.Add(name, strconv.FormatFloat(field.Float(), 'g', -1, 64))

		case reflect.Bool:
			if field.Bool() {
				params.Add(name, "true")
			} else {
				params.Add(name, "false")
			}

		case reflect.Slice:
			for i := 0; i < field.Len(); i++ {
				elem := field.Index(i)
				switch elem.Kind() {
				case reflect.String:
					params.Add(name, elem.String())
				case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
					params.Add(name, strconv.FormatInt(elem.Int(), 10))
				case reflect.Bool:
					if elem.Bool() {
						params.Add(name, "true")
					} else {
						params.Add(name, "false")
					}
				default:
					return "", fmt.Errorf("неподдерживаемый тип элемента среза: %s", elem.Kind())
				}
			}
		}

	}
	return params.Encode(), nil

}
