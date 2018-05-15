package get

import (
	"errors"
	"fmt"
	"github.com/dima-kov/go-rarser/vars"
	"net/http"
	"reflect"
)

type coreGetParser struct{}

func NewCoreGetParser() Parser {
	return coreGetParser{}
}

func (p coreGetParser) ParseGET(r *http.Request, structField *reflect.StructField) (string, error) {
	paramValue := r.URL.Query().Get(structField.Tag.Get(vars.TagGet))
	if paramValue != "" {
		return paramValue, nil
	}
	defaultValue, ok := structField.Tag.Lookup(vars.TagDefault)
	if !ok {
		return "", errors.New(
			fmt.Sprintf("empty required GET param: %s", structField.Tag.Get(vars.TagGet)),
		)
	}
	return defaultValue, nil
}
