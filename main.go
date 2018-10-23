package rarser

import (
	"github.com/dima-kov/go-rarser/body"
	"github.com/dima-kov/go-rarser/get"
	"github.com/dima-kov/go-rarser/path"
	"net/http"
)

var instance RequestParser

func init() {
	instance = newRequestParser(defaultPath, defaultBody, defaultGet)
}

func Parse(r *http.Request, parseInto interface{}) {
	instance.Parse(r, parseInto)
}

func SetPathParser(parser path.Parser) {
	instance.setPathParser(parser)
}

func SetBodyParser(parser body.Parser) {
	instance.setBodyParser(parser)
}

func SetGetParser(parser get.Parser) {
	instance.setGetParser(parser)
}
