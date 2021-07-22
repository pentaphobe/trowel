package trowel

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

type ErrorsSuite struct {
	suite.Suite
}

func (t *ErrorsSuite) TestTrowelError() {
	msg := "MESSAGE"
	err := NewError(nil, msg)
	t.Contains(err.Message, msg)
	t.Contains(err.Error(), "TrowelError: "+msg)
}

func (t *ErrorsSuite) TestTrowelIndexError() {
	msg := "MESSAGE"
	err := NewIndexError(nil, msg+" %d", 31337)
	t.Contains(err.Message, msg+" 31337")
	t.Contains(err.Error(), "TrowelIndexError: "+msg+" 31337")
}

func (t *ErrorsSuite) TestTrowelKeyError() {
	msg := "MESSAGE"
	err := NewKeyError(nil, msg+" %s", "foo")
	t.Contains(err.Message, msg+" foo")
	t.Contains(err.Error(), "TrowelKeyError: "+msg+" foo")
}

func (t *ErrorsSuite) TestTrowelParseError() {
	msg := "MESSAGE"
	err := NewParseError(nil, msg)
	t.Contains(err.Message, msg)
	t.Contains(err.Error(), "TrowelParseError: MESSAGE")
}

func (t *ErrorsSuite) TestErrorHelpers() {
	// intentionally skip constructor, resulting in empty `errors` array
	obj := &trowelWrapper{}
	t.Nil(obj.Error())
	obj.addError(NewError(obj, "hi"))
	t.Len(obj.Errors(), 1)
	t.NotNil(obj.Error())
}

func TestErrorsSuite(t *testing.T) {
	suite.Run(t, new(ErrorsSuite))
}
