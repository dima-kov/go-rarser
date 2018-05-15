# Go-Rarser

Go library to parse all http params into one struct

## Install

    go get github.com/dima-kov/go-rarser


## Usage

1. Define struct to which parse all values

    ```
	type SaveUserParams struct {
		ID uint `path:"id"`
		From time.Time `get:"from"`
		Body struct {
			Name string `json:"name"`
		} `body:"true"`
	}
	```

1. Use parser in handler to parse them all:
    ```
	params := SaveUserParams{}
	err := rarser.GetRequestParser().Parse(r, &params)
	if err != nil {
		// check errors
		w.Write([]byte(err.Error()))
		return
	}
	```
1. Use value:
    ```
	fmt.Fprintf("%v, %v, %v", params.ID, params.From, params.Body.Name)
	```

By default `goji.io` is used as http router to parse values from URL path, to override it read next section.

## Extended usage

Call `rarser.InitRequestParser(gorouterParser.NewGoRouterPathParser())` method on the application start, for example, in package with routing:

```
import (
	"net/http"
	"github.com/dima-kov/go-rarser"
	gorouterParser "github.com/dima-kov/go-rarser/path/gorouter"
	"log"
	"github.com/vardius/gorouter"
)

func main() {
	userHandler := handlers.NewUserHandler()
	rarser.InitRequestParser(gorouterParser.NewGoRouterPathParser())

	router := gorouter.New()
	router.POST("/user/{id}", http.HandlerFunc(CreateUser))

	log.Fatal(http.ListenAndServe(":8000", router))
}

func (h *UserHandlers) CreateUser(w http.ResponseWriter, r *http.Request) {
	params := SaveUserParams{}
	err := rarser.GetRequestParser().Parse(r, &params)
	if err != nil {
		// check errors
		w.Write([]byte(err.Error()))
		return
	}
	fmt.Fprintf("%v, %v, %v", params.ID, params.From, params.Body.Name)
}
```