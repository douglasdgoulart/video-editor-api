package validator

import (
	"fmt"
	"reflect"
	"strings"
)

func ValidateRequiredFields(v interface{}) error {
	return validateRequiredFields(v, "")
}

func validateRequiredFields(v interface{}, parentPath string) error {
	val := reflect.ValueOf(v)
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}

	typ := val.Type()
	for i := 0; i < val.NumField(); i++ {
		field := val.Field(i)
		fieldType := typ.Field(i)
		jsonTag := getJSONTag(fieldType)
		fieldPath := buildFieldPath(parentPath, jsonTag)

		if fieldType.Tag.Get("required") == "true" {
			if isEmptyValue(field) {
				return fmt.Errorf("field %s is required", fieldPath)
			}
		}

		if field.Kind() == reflect.Struct {
			if field.CanAddr() {
				if err := validateRequiredFields(field.Addr().Interface(), fieldPath); err != nil {
					return err
				}
			} else {
				if err := validateRequiredFields(field.Interface(), fieldPath); err != nil {
					return err
				}
			}
		}
	}
	return nil
}

func buildFieldPath(parentPath, fieldName string) string {
	if parentPath == "" {
		return fieldName
	}
	return fmt.Sprintf("%s.%s", parentPath, fieldName)
}

func getJSONTag(field reflect.StructField) string {
	tag := field.Tag.Get("json")
	if tag == "" || tag == "-" {
		return field.Name
	}
	return strings.Split(tag, ",")[0]
}

func isEmptyValue(v reflect.Value) bool {
	switch v.Kind() {
	case reflect.String, reflect.Array, reflect.Map, reflect.Slice:
		return v.Len() == 0
	case reflect.Bool:
		return !v.Bool()
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return v.Int() == 0
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return v.Uint() == 0
	case reflect.Float32, reflect.Float64:
		return v.Float() == 0
	case reflect.Interface, reflect.Ptr:
		return v.IsNil()
	}
	return false
}
