package get

import (
	"net/http"
	"reflect"
)

type Parser interface {
	ParseGET(r *http.Request, structField *reflect.StructField) (string, error)
}
