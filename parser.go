package rarser

import (
	"github.com/dima-kov/go-rarser/body"
	"github.com/dima-kov/go-rarser/get"
	"github.com/dima-kov/go-rarser/path"
	"github.com/dima-kov/go-rarser/vars"
	"net/http"
	"reflect"
	"strconv"
	"strings"
	"time"
)

type RequestParser interface {
	Parse(r *http.Request, parseTo interface{}) error
}

type requestParser struct {
	pathParser path.Parser
	getParser  get.Parser
	bodyParser body.Parser
}

type structValueField struct {
	valueField  *reflect.Value
	structField *reflect.StructField
}

func (p requestParser) Parse(r *http.Request, parseInto interface{}) error {
	parseIntoFields := p.getFields(reflect.ValueOf(parseInto).Elem())
	for _, field := range parseIntoFields {
		if err := p.parseStructField(r, field.valueField, field.structField); err != nil {
			return err
		}
	}
	return nil
}

// Returns fields of struct passed as (reflect.ValueOf(structObj))
// For embedding structs: recursively append to list
func (p requestParser) getFields(reflectField reflect.Value) []structValueField {
	fields := make([]structValueField, 0)

	for i := 0; i < reflectField.NumField(); i++ {
		valueField := reflectField.Field(i)
		structField := reflectField.Type().Field(i)
		if valueField.Kind() == reflect.Struct && structField.Anonymous {
			fields = append(fields, p.getFields(valueField)...)
		} else {
			fields = append(fields, structValueField{
				&valueField,
				&structField,
			})
		}
	}
	return fields
}

func getTagType(tag reflect.StructTag) string {
	return strings.Split(string(tag), ":")[0]
}

func (p requestParser) parseStructField(r *http.Request, field *reflect.Value, structField *reflect.StructField) error {
	switch getTagType(structField.Tag) {
	case vars.TagPath:
		err := p.parseFieldFromUrl(r, field, structField)
		if err != nil {
			return err
		}
	case vars.TagGet:
		err := p.parseFieldFromGet(r, field, structField)
		if err != nil {
			return err
		}
	case vars.TagBody:
		err := p.parseBody(r, field)
		if err != nil {
			return err
		}
	}
	return nil
}

func (p requestParser) parseFieldFromUrl(r *http.Request, field *reflect.Value, structField *reflect.StructField) error {
	return p.setValue(field, p.pathParser.ParsePath(r, structField))
}

func (p requestParser) parseFieldFromGet(r *http.Request, field *reflect.Value, structField *reflect.StructField) error {
	value, err := p.getParser.ParseGET(r, structField)
	if err != nil {
		return err
	}
	return p.setValue(field, value)
}

func (p requestParser) parseBody(r *http.Request, field *reflect.Value) error {
	return p.bodyParser.ParseBody(r, field)

}

// Sets value to field converting value to field type
func (p requestParser) setValue(field *reflect.Value, value string) error {
	switch field.Kind() {
	case reflect.String:
		field.SetString(value)
	case reflect.Int:
		return p.convertAndSetInt(value, field)
	case reflect.Uint:
		return p.convertAndSetUint(value, field)
	case reflect.Float64:
		return p.convertAndSetFloat(value, field)
	case reflect.TypeOf(time.Time{}).Kind():
		return p.convertAndSetTime(value, field)
	case reflect.Bool:
		return p.convertAndSetBool(value, field)
	}
	return nil
}

func (requestParser) convertAndSetInt(value string, field *reflect.Value) error {
	converted, err := strconv.Atoi(value)
	if err != nil {
		return err
	}
	field.SetInt(int64(converted))
	return nil
}

func (requestParser) convertAndSetUint(value string, field *reflect.Value) error {
	converted, err := strconv.Atoi(value)
	if err != nil {
		return err
	}
	field.SetUint(uint64(converted))
	return nil
}

func (p requestParser) convertAndSetFloat(value string, field *reflect.Value) error {
	converted, err := strconv.ParseFloat(value, 64)
	if err != nil {
		return err
	}
	field.SetFloat(converted)
	return nil
}

func (p requestParser) convertAndSetTime(value string, field *reflect.Value) error {
	converted, err := time.Parse(time.RFC3339, value)
	if err != nil {
		return err
	}
	field.Set(reflect.ValueOf(converted.UTC()))
	return nil
}

func (p requestParser) convertAndSetBool(value string, field *reflect.Value) error {
	result := false
	if value == "true" {
		result = true
	}
	field.SetBool(result)
	return nil
}
