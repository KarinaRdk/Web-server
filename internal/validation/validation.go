package validation

import (
	"fmt"
	"reflect"
	"time"
)

func Check(i interface{}) bool {
	val := reflect.ValueOf(i)
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}

	if val.Kind() != reflect.Struct {
		fmt.Println("Error: validate function expects a struct or a pointer to a struct.")
		return false
	}

	return checkStruct(val)
}
func checkStruct(val reflect.Value) bool {
	for i := 0; i < val.NumField(); i++ {
		field := val.Field(i)
		fieldType := val.Type().Field(i)

		// Skip checking for nil on *time.Time fields
		if fieldType.Type.Kind() == reflect.Ptr && fieldType.Type.Elem().Kind() == reflect.Struct && fieldType.Type.Elem() == reflect.TypeOf(time.Time{}) {
			continue
		}

		switch field.Kind() {
		case reflect.Ptr:
			if field.IsNil() {
				fmt.Printf("Field %s is nil\n", fieldType.Name)
				return false
			}
			if field.Elem().Kind() == reflect.Struct {
				if !checkStruct(field.Elem()) {
					return false
				}
			}
		case reflect.Slice:
			if field.IsNil() {
				fmt.Printf("Field %s is nil\n", fieldType.Name)
				return false
			}
			for j := 0; j < field.Len(); j++ {
				if !checkStruct(field.Index(j)) {
					return false
				}
			}
		case reflect.Struct:
			if !checkStruct(field) {
				return false
			}
		}
	}

	return true
}
