package trowel

import (
	"fmt"
)

type TrowelError struct {
	Message string
}

type TrowelIndexError struct {
	TrowelError
	Index int
}

type TrowelKeyError struct {
	TrowelError
	Key string
}

type TrowelParseError struct {
	TrowelError
}

func (e *TrowelError) Error() string {
	return fmt.Sprintf("TrowelError: %s", e.Message)
}

func (e *TrowelIndexError) Error() string {
	return fmt.Sprintf("TrowelIndexError: %s", e.Message)
}

func (e *TrowelKeyError) Error() string {
	return fmt.Sprintf("TrowelKeyError: %s", e.Message)
}

func (e *TrowelParseError) Error() string {
	return fmt.Sprintf("TrowelParseError: %s", e.Message)
}

func NewError(message string, rest ...interface{}) *TrowelError {
	return &TrowelError{
		Message: fmt.Sprintf(message, rest...),
	}
}

func NewIndexError(message string, index int, rest ...interface{}) *TrowelIndexError {
	return &TrowelIndexError{
		TrowelError: TrowelError{
			Message: fmt.Sprintf(message, append([]interface{}{index}, rest...)...),
		},
		Index: index,
	}
}

func NewKeyError(message string, key string, rest ...interface{}) *TrowelKeyError {
	return &TrowelKeyError{
		TrowelError: TrowelError{
			Message: fmt.Sprintf(message, append([]interface{}{key}, rest...)...),
		},
		Key: key,
	}
}

func NewParseError(message string, rest ...interface{}) *TrowelParseError {
	return &TrowelParseError{
		TrowelError: TrowelError{
			Message: fmt.Sprintf(message, rest...),
		},
	}
}
