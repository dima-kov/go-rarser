package rarser

import (
	"github.com/dima-kov/go-rarser/body"
	"github.com/dima-kov/go-rarser/get"
	"github.com/dima-kov/go-rarser/path/gorouter"
	"sync"
)

var instance *requestParser
var once sync.Once

func GetRequestParser() RequestParser {
	once.Do(func() {
		pathParser := gorouter.NewGoRouterPathParser()
		getParser := get.NewCoreGetParser()
		bodyParser := body.NewJsonBodyParser()
		instance = &requestParser{
			pathParser,
			getParser,
			bodyParser,
		}
	})
	return instance
}
