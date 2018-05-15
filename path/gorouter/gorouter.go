package gorouter

import (
	"github.com/dima-kov/go-rarser/path"
	"github.com/dima-kov/go-rarser/vars"
	"github.com/vardius/gorouter"
	"net/http"
	"reflect"
)

type goRouterPathParser struct{}

func NewGoRouterPathParser() path.Parser {
	return goRouterPathParser{}
}

func (p goRouterPathParser) ParsePath(r *http.Request, structField *reflect.StructField) string {
	params, _ := gorouter.FromContext(r.Context())
	return params.Value(structField.Tag.Get(vars.TagPath))
}
