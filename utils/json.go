package utils

import (
	"reflect"
	"strings"
)

func GetJSONFieldName(obj interface{}, fieldName string) string {
	objType := reflect.TypeOf(obj)
	if objType.Kind() == reflect.Ptr {
		objType = objType.Elem()
	}

	for i := 0; i < objType.NumField(); i++ {
		field := objType.Field(i)
		jsonTag := field.Tag.Get("json")
		if field.Name == fieldName {
			if jsonTag != "" {
				return strings.Split(jsonTag, ",")[0]
			}
			return fieldName
		}
	}

	return fieldName
}
