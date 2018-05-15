package body

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"reflect"
)

type jsonBodyParser struct{}

func NewJsonBodyParser() Parser {
	return jsonBodyParser{}
}

func (p jsonBodyParser) ParseBody(r *http.Request, field *reflect.Value) error {
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
