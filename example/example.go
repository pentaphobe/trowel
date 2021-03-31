package main

import (
	"encoding/json"

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
	var obj interface{}
	err := json.Unmarshal([]byte(exampleJSON), &obj)
	if err != nil {
		panic(err)
	}

	// Wrap the object
	t := trowel.NewTrowel(obj)

	result, err := t.Path(`.foo.baz[3].bling`)
	if err != nil {
		panic(err)
	}
	println(result.Get().(bool))

	result, err = t.Path(`."key.with.special_characters"`)
	if err != nil {
		panic(err)
	}
	println(result.Get().(bool))
}
