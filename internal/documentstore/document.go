package documentstore

import (
	"fmt"
	"reflect"
)

type DocumentFieldType string

const (
	DocumentFieldTypeString DocumentFieldType = "string"
	DocumentFieldTypeNumber DocumentFieldType = "int"
	DocumentFieldTypeBool   DocumentFieldType = "bool"
	DocumentFieldTypeArray  DocumentFieldType = "array"
	DocumentFieldTypeObject DocumentFieldType = "object"
)

type DocumentField struct {
	Type  DocumentFieldType `json:"type"`
	Value any               `json:"value"`
}

type Document struct {
	Fields map[string]DocumentField `json:"fields"`
}

type MyStruct struct {
	X int
}

func MarshalDocument(input interface{}) (*Document, error) {
	v := reflect.ValueOf(input)
	t := reflect.TypeOf(input)

	if t.Kind() == reflect.Ptr {
		v = v.Elem()
		t = t.Elem()
	}

	if t.Kind() != reflect.Struct {
		return nil, nil
	}
	doc := Document{
		Fields: make(map[string]DocumentField),
	}

	for i := 0; i < v.NumField(); i++ {
		f := t.Field(i)
		val := v.Field(i)

		var FieldType DocumentFieldType
		var FieldValue interface{}

		switch val.Kind() {
		case reflect.String:
			FieldType = DocumentFieldTypeString
			FieldValue = val.String()
		case reflect.Int:
			FieldType = DocumentFieldTypeNumber
			FieldValue = val.Int()
		case reflect.Bool:
			FieldType = DocumentFieldTypeBool
			FieldValue = val.Bool()
		case reflect.Slice:
			FieldType = DocumentFieldTypeArray
			FieldValue = val.Interface()
		case reflect.Struct:
			FieldType = DocumentFieldTypeObject
			FieldValue = val.Interface()
		default:
			continue
		}
		doc.Fields[f.Name] = DocumentField{
			Type:  FieldType,
			Value: FieldValue,
		}
	}
	return &doc, nil
}

func UnmarshalDocument(doc *Document, output any) error {
	v := reflect.ValueOf(output)

	if v.Kind() != reflect.Ptr || v.Elem().Kind() != reflect.Struct {
		return fmt.Errorf("output is not a pointer to a struct")
	}
	stValue := v.Elem()
	stType := stValue.Type()
	for i := 0; i < stType.NumField(); i++ {
		f := stType.Field(i)
		fValue := stValue.Field(i)
		if !fValue.CanSet() {
			continue
		}
		jsonTag := f.Tag.Get("json")
		if jsonTag == "" {
			jsonTag = f.Name
		}
		if val, ok := doc.Fields[jsonTag]; ok {
			valR := reflect.ValueOf(val.Value)
			if valR.Type().AssignableTo(fValue.Type()) {
				fValue.Set(valR)
			} else {
				return fmt.Errorf("type mismatch for field %s: expected %s but got %s", f.Name, fValue.Type(), valR.Type())
			}
		}
	}
	return nil
}
