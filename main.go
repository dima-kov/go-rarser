package rarser

import (
	"github.com/dima-kov/go-rarser/body"
	"github.com/dima-kov/go-rarser/get"
	"github.com/dima-kov/go-rarser/path"
	"github.com/dima-kov/go-rarser/path/goji"
	"sync"
)

var instance *requestParser
var once sync.Once

func InitRequestParser(parsers ...interface{}) RequestParser {
	once.Do(func() {
		pathParser := goji.NewGojiPathParser()
		getParser := get.NewCoreGetParser()
		bodyParser := body.NewJsonBodyParser()
		if len(parsers) > 0 && parsers[0] != nil {
			pathParser = parsers[0].(path.Parser)
		}
		if len(parsers) > 1 && parsers[1] != nil {
			getParser = parsers[1].(get.Parser)
		}
		if len(parsers) > 2 && parsers[2] != nil {
			bodyParser = parsers[2].(body.Parser)
		}
		instance = &requestParser{
			pathParser,
			getParser,
			bodyParser,
		}
	})
	return instance
}

func GetRequestParser() RequestParser {
	if instance == nil {
		return InitRequestParser()
	}
	return instance
}
