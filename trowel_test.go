package trowel

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/suite"
)

type TrowelSuite struct {
	suite.Suite
}

func (suite *TrowelSuite) TestNew() {
	_ = NewTrowel(nil)
}

func (suite *TrowelSuite) TestGet() {
	data := "hi"
	t := NewTrowel(data)
	suite.Equal(data, t.Get())
}

func (suite *TrowelSuite) TestArray() {
	data := []interface{}{}
	t := NewTrowel(data)
	arr, err := t.Array()
	suite.Nil(err)
	suite.Equal(arr, data)
}

func (suite *TrowelSuite) TestMap() {
	data := make(map[string]interface{})
	t := NewTrowel(data)
	arr, err := t.Map()
	suite.Nil(err)
	suite.Equal(arr, data)
}

func (suite *TrowelSuite) TestIndex() {
	data := []interface{}{
		"hi",
		1,
	}
	t := NewTrowel(data)
	for i, val := range data {
		result, err := t.Index(i)
		suite.Nil(err)
		suite.Equal(val, result.Get())
	}
}

func (suite *TrowelSuite) TestKey() {
	data := map[string]interface{}{
		"hi":    "hello",
		"howdy": 1,
	}
	t := NewTrowel(data)
	for key, val := range data {
		result, err := t.Key(key)
		suite.Nil(err)
		suite.Equal(val, result.Get())
	}
}

func (suite *TrowelSuite) TestDeep() {
	input := `
	{
		"foo": {
			"bar": 1,
			"baz": [1, "boffle", false, {
				"bling": true
			}],
			"bing": {
				"bof": true,
				"bif": 1.3
			}
		}
	}
	`
	var data interface{}
	err := json.Unmarshal([]byte(input), &data)
	suite.Nil(err)
	suite.NotNil(data)
	t := NewTrowel(data)

	first, err := t.Key("foo")
	suite.Nil(err)
	arr, err := first.Key("baz")
	suite.Nil(err)
	arrMap, err := arr.Index(3)
	suite.Nil(err)
	arrMapKey, err := arrMap.Key("bling")
	suite.Nil(err)
	suite.Equal(arrMapKey.Get().(bool), true)
}

func (suite *TrowelSuite) TestPath() {
	input := `
	{
		"foo": {
			"bar": 1,
			"baz": [1, "boffle", false, {
				"bling": true
			}],
			"bing": {
				"bof": true,
				"bif": 1.3
			}
		}
	}
	`
	var data interface{}
	err := json.Unmarshal([]byte(input), &data)
	suite.Nil(err)
	suite.NotNil(data)
	t := NewTrowel(data)

	paths := map[string]interface{}{
		// No ints in JSON, so have to assume float64
		"foo.bar":          1.0,
		"foo.baz[0]":       1.0,
		"foo.baz[1]":       "boffle",
		"foo.baz[3].bling": true,
	}
	for path, expected := range paths {
		child, err := t.Path(path)
		suite.Nil(err)
		suite.Equal(child.Get(), expected)
	}
}

func TestTrowelSuite(t *testing.T) {
	suite.Run(t, new(TrowelSuite))
}
