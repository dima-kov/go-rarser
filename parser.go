package rarser

import (
	"reflect"
	"sync"
	"goji.io/pat"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"time"
)

const TagUrl = "url"
const TagBody = "body"
const TagGet = "get"
const TagDefault = "default"

type RequestParser interface {
	Parse(r *http.Request, parseTo interface{}) error
}

type requestParser struct{}

var instance *requestParser
var once sync.Once

type structValueField struct {
	valueField  *reflect.Value
	structField *reflect.StructField
}

func (p requestParser) Parse(r *http.Request, parseTo interface{}) error {
	fields := p.getFields(reflect.ValueOf(parseTo).Elem())
	for _, field := range fields {
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
	case TagUrl:
		err := p.parseFieldFromUrl(r, field, structField)
		if err != nil {
			return err
		}
	case TagGet:
		err := p.parseFieldFromGet(r, field, structField)
		if err != nil {
			return err
		}
	case TagBody:
		err := p.parseBody(r, field)
		if err != nil {
			return err
		}
	}
	return nil
}

func (p requestParser) parseFieldFromUrl(r *http.Request, field *reflect.Value, structField *reflect.StructField) error {
	paramValue := pat.Param(r, structField.Tag.Get(TagUrl))
	return p.setValue(field, paramValue)
}

func (p requestParser) parseFieldFromGet(r *http.Request, field *reflect.Value, structField *reflect.StructField) error {
	paramValue := r.URL.Query().Get(structField.Tag.Get(TagGet))
	if paramValue != "" {
		return p.setValue(field, paramValue)
	}
	defaultValue, ok := structField.Tag.Lookup(TagDefault)
	if !ok {
		return errors.New(
			fmt.Sprintf("empty required GET param: %s", structField.Tag.Get(TagGet)),
		)
	}
	return p.setValue(field, defaultValue)
}

func (p requestParser) parseBody(r *http.Request, field *reflect.Value) error {
	b, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		return err
	}
	newValue := reflect.New(field.Type())
	err = json.Unmarshal(b, newValue.Interface())
	field.Set(newValue.Elem())
	return err
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
