package main

import (
	"encoding/json"
	"fmt"

	"github.com/pentaphobe/trowel"
)

const exampleJSON = `{
  "foo": {
    "baz": [1, "boffle", false, {
      "bling": true
    }]
  },
	"key.with.special_characters": true
}
`

func main() {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("recovered:", r)
		}
	}()
	var obj interface{}
	err := json.Unmarshal([]byte(exampleJSON), &obj)
	if err != nil {
		panic(err)
	}

	// Wrap the object in a Trowel cursor
	t := trowel.NewTrowel(obj)

	result := t.Path(`.foo.baz[3].bling`)
	if result.HasErrors() {
		panic(result.Error())
	}
	println(result.Get().(bool))

	// equivalent to above
	result = t.Key("foo").Key("baz").Index(3).Key("bling")
	if result.HasErrors() {
		panic(result.Error())
	}
	println(result.Get().(bool))

	// quoted keys
	result = t.Path(`."key.with.special_characters"`)
	if result.HasErrors() {
		panic(result.Error())
	}
	println(result.Get().(bool))

	// invalid access also returns a cursor, so chaining is safe
	result = t.Key("foo").Key("NO-OP").Index(3000).Key("bling")
	if result.HasErrors() {
		panic(result.Error())
	}
	println(result.Get().(bool))
}
