package trowel

import (
	"fmt"
)

type ErrorContext interface{}
type TrowelError struct {
	Context ErrorContext
	Message string
}

type TrowelIndexError struct {
	Context ErrorContext
	Message string
	Index   int
}

type TrowelKeyError struct {
	Context ErrorContext
	Message string
	Key     string
}

type TrowelParseError struct {
	Context ErrorContext
	Message string
}

func (e *TrowelError) Error() string {
	return fmt.Sprintf("TrowelError: %s, CONTEXT: %+v", e.Message, e.Context)
}

func (e *TrowelIndexError) Error() string {
	return fmt.Sprintf("TrowelIndexError: %s, CONTEXT: %+v", e.Message, e.Context)
}

func (e *TrowelKeyError) Error() string {
	return fmt.Sprintf("TrowelKeyError: %s, CONTEXT: %+v", e.Message, e.Context)
}

func (e *TrowelParseError) Error() string {
	return fmt.Sprintf("TrowelParseError: %s, CONTEXT: %+v", e.Message, e.Context)
}

func NewError(ctx ErrorContext, message string, rest ...interface{}) *TrowelError {
	return &TrowelError{
		Context: ctx,
		Message: fmt.Sprintf(message, rest...),
	}
}

func NewIndexError(ctx ErrorContext, message string, index int, rest ...interface{}) *TrowelIndexError {
	return &TrowelIndexError{
		Context: ctx,
		Message: fmt.Sprintf(message, append([]interface{}{index}, rest...)...),
		Index:   index,
	}
}

func NewKeyError(ctx ErrorContext, message string, key string, rest ...interface{}) *TrowelKeyError {
	return &TrowelKeyError{
		Context: ctx,
		Message: fmt.Sprintf(message, append([]interface{}{key}, rest...)...),
		Key:     key,
	}
}

func NewParseError(ctx ErrorContext, message string, rest ...interface{}) *TrowelParseError {
	return &TrowelParseError{
		Context: ctx,
		Message: fmt.Sprintf(message, rest...),
	}
}
