package rarser

import (
	"github.com/dima-kov/go-rarser/body"
	"github.com/dima-kov/go-rarser/get"
	"github.com/dima-kov/go-rarser/path/goji"
)

var defaultBody = body.NewJsonBodyParser()
var defaultGet = get.NewCoreGetParser()
var defaultPath = goji.NewGojiPathParser()
