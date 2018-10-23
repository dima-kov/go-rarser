package rarser

import (
	"fmt"
	"testing"
)

func TestParserCreated(t *testing.T) {
	parser, ok := instance.(*requestParser)
	if ok {
		fmt.Println("SUPER")
	}
	fmt.Println(parser.bodyParser)
	fmt.Println(parser.pathParser)
	fmt.Println(parser.getParser)
}
