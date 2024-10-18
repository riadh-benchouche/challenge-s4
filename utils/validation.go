package utils

import (
	"github.com/go-playground/validator/v10"
	"mime/multipart"
	"reflect"
	"strings"
)

func GetValidationErrors(errors validator.ValidationErrors, jsonBody interface{}) map[string]string {
	var validationErrors = make(map[string]string)

	for _, e := range errors {
		fieldParts := strings.Split(e.StructNamespace(), ".")
		fieldPartsWithoutNamespace := fieldParts
		if len(fieldParts) > 1 {
			fieldPartsWithoutNamespace = fieldParts[1:]
		}
		currentBody := jsonBody
		fieldNameJSON := GetJSONFieldName(currentBody, fieldPartsWithoutNamespace[0])

		for i := 0; i < len(fieldPartsWithoutNamespace)-1; i++ {
			currentBody = getField(currentBody, fieldPartsWithoutNamespace[i])
			if currentBody == nil {
				break
			}

			if currentBody != nil {
				fieldNameJSON += "." + GetJSONFieldName(currentBody, fieldPartsWithoutNamespace[len(fieldPartsWithoutNamespace)-1])
			}
		}

		if fieldNameJSON != "" {
			validationErrors[fieldNameJSON] = e.Tag()
		}
	}

	return validationErrors
}

func getField(jsonBody interface{}, fieldName string) interface{} {
	v := reflect.ValueOf(jsonBody)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	if v.Kind() != reflect.Struct {
		return nil
	}

	field := v.FieldByName(fieldName)
	if !field.IsValid() {
		return nil
	}

	return field.Interface()
}

// image formats and magic numbers
var magicTable = map[string]string{
	"\xff\xd8\xff":      "image/jpeg",
	"\x89PNG\r\n\x1a\n": "image/png",
	"GIF87a":            "image/gif",
	"GIF89a":            "image/gif",
}

func MimeFromIncipit(incipit []byte) string {
	incipitStr := string(incipit)
	for magic, mime := range magicTable {
		if strings.HasPrefix(incipitStr, magic) {
			return mime
		}
	}

	return ""
}

func IsImage(file multipart.FileHeader) bool {
	incipit := make([]byte, 512)
	fileSrc, err := file.Open()
	if err != nil {
		return false
	}
	defer fileSrc.Close()

	_, err = fileSrc.Read(incipit)
	if err != nil {
		return false
	}

	return MimeFromIncipit(incipit) != ""
}

func GetImageExt(file multipart.FileHeader) string {
	incipit := make([]byte, 512)
	fileSrc, err := file.Open()
	if err != nil {
		return ""
	}
	defer fileSrc.Close()

	_, err = fileSrc.Read(incipit)
	if err != nil {
		return ""
	}

	mime := MimeFromIncipit(incipit)
	if mime != "" {
		switch mime {
		case "image/jpeg":
			return ".jpg"
		case "image/png":
			return ".png"
		case "image/gif":
			return ".gif"
		}
	}

	return ""
}
