package trowel

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/stretchr/testify/suite"
)

type TrowelSuite struct {
	suite.Suite
}

func (suite *TrowelSuite) TestNew() {
	_ = NewTrowel(nil)
}

func (suite *TrowelSuite) TestParsePath() {
	data := map[string][]interface{}{
		".key":                              []interface{}{"key"},
		".many.keys":                        []interface{}{"many", "keys"},
		"[0]":                               []interface{}{0},
		"[1][2]":                            []interface{}{1, 2},
		".complex[1].entry.with[1][2].lots": []interface{}{"complex", 1, "entry", "with", 1, 2, "lots"},
	}
	for input, expected := range data {
		result, err := parsePath(input)
		suite.Nil(err)
		suite.Equal(expected, result, fmt.Sprintf("input:'%s'", input))
	}
}

func (suite *TrowelSuite) TestParsePathInvalid() {
	data := map[string]interface{}{
		"nodot": &TrowelParseError{},
		"[NaN]": &TrowelParseError{},
	}
	for input, expected := range data {
		result, err := parsePath(input)
		suite.Nil(result)
		suite.IsType(expected, err, fmt.Sprintf("input:'%s'", input))
	}
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
	suite.NotNil(arr)
	suite.Equal(arr, data)
}

func (suite *TrowelSuite) TestArrayInvalid() {
	data := []struct {
		value    interface{}
		expected interface{}
	}{
		{nil, &TrowelError{}},
		{"not an array", &TrowelError{}},
	}

	for _, entry := range data {
		t := NewTrowel(entry.value)
		result, err := t.Array()
		suite.Nil(result)
		suite.IsType(entry.expected, err)
	}
}

func (suite *TrowelSuite) TestMap() {
	data := make(map[string]interface{})
	t := NewTrowel(data)
	arr, err := t.Map()
	suite.Nil(err)
	suite.NotNil(arr)
	suite.Equal(arr, data)
}

func (suite *TrowelSuite) TestMapInvalid() {
	data := []struct {
		value    interface{}
		expected interface{}
	}{
		{nil, &TrowelError{}},
		{"not a map", &TrowelError{}},
	}

	for _, entry := range data {
		t := NewTrowel(entry.value)
		result, err := t.Map()
		suite.Nil(result)
		suite.IsType(entry.expected, err)
	}
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
		suite.NotNil(val)
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
		suite.NotNil(val)
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
			},
			"special.key": true
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
		".foo.bar":           1.0,
		".foo.baz[0]":        1.0,
		".foo.baz[1]":        "boffle",
		".foo.baz[3].bling":  true,
		`.foo."special.key"`: true,
	}
	for path, expected := range paths {
		child, err := t.Path(path)
		suite.Nil(err)
		suite.NotNil(child)
		suite.Equal(expected, child.Get())
	}
}

func (suite *TrowelSuite) TestPath_Invalid() {
	input := `
	{
		"foo": {
			"bar": 1,
			"baz": [1, "boffle", false, {
				"bling": true
			}],
			"nullkey": null,
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
		".nonexistent":           &TrowelKeyError{},
		".foo.bar.not-an-object": &TrowelKeyError{},
		".foo.nullkey.invalid":   &TrowelKeyError{},
		".foo.baz[invalid]":      &TrowelParseError{},
	}
	for path, expected := range paths {
		val, err := t.Path(path)
		suite.Nil(val)
		suite.IsType(expected, err, fmt.Sprintf("path:'%s'", path))
	}
}

func (suite *TrowelSuite) TestIndex_Nil() {
	t := NewTrowel(nil)

	result, err := t.Index(0)
	suite.Nil(result.Get())
	suite.IsType(err, &TrowelIndexError{})
}

func (suite *TrowelSuite) TestIndex_Non_Array() {
	t := NewTrowel(10)

	result, err := t.Index(0)
	suite.Nil(result.Get())
	suite.IsType(err, &TrowelIndexError{})
}

func TestTrowelSuite(t *testing.T) {
	suite.Run(t, new(TrowelSuite))
}
