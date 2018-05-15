package goji

import (
	"github.com/dima-kov/go-rarser/vars"
	"github.com/dima-kov/go-rarser/path"
	"goji.io/pat"
	"net/http"
	"reflect"
)

type gojiPathParser struct{}

func NewGojiPathParser() path.Parser {
	return gojiPathParser{}
}

func (p gojiPathParser) ParsePath(r *http.Request, structField *reflect.StructField) string {
	return pat.Param(r, structField.Tag.Get(vars.TagPath))
}
