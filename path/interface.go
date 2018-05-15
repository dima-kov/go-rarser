package path

import (
	"net/http"
	"reflect"
)

type Parser interface {
	ParsePath(r *http.Request, structField *reflect.StructField) string
}
