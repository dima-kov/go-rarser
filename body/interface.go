package body

import (
	"net/http"
	"reflect"
)

type Parser interface {
	ParseBody(r *http.Request, field *reflect.Value) error
}
