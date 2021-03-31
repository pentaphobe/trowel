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
	err := NewError(msg)
	t.Equal(msg, err.Message)
	t.Equal("TrowelError: "+msg, err.Error())
}

func (t *ErrorsSuite) TestTrowelIndexError() {
	msg := "MESSAGE"
	err := NewIndexError(msg+" %d", 31337)
	t.Equal(msg+" 31337", err.TrowelError.Message)
	t.Equal("TrowelIndexError: "+msg+" 31337", err.Error())
}

func (t *ErrorsSuite) TestTrowelKeyError() {
	msg := "MESSAGE"
	err := NewKeyError(msg+" %s", "foo")
	t.Equal(msg+" foo", err.TrowelError.Message)
	t.Equal("TrowelKeyError: "+msg+" foo", err.Error())
}

func (t *ErrorsSuite) TestTrowelParseError() {
	msg := "MESSAGE"
	err := NewParseError(msg)
	t.Equal(msg, err.TrowelError.Message)
	t.Equal("TrowelParseError: MESSAGE", err.Error())
}

func TestErrorsSuite(t *testing.T) {
	suite.Run(t, new(ErrorsSuite))
}
