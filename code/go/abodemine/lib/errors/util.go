package errors

import (
	std_errors "errors"
)

var As = std_errors.As
var Is = std_errors.Is
var New = std_errors.New

func Callback(err error, f func(error)) {
	if err != nil {
		f(err)
	}
}

// First returns the first error in the chain.
func First(err error) *Object {
	switch v := err.(type) {
	case *Object:
		return v
	case *Chain:
		return v.First()
	}

	return new(Chain).First()
}

// Last returns the last error in the chain.
func Last(err error) *Object {
	switch v := err.(type) {
	case *Object:
		return v
	case *Chain:
		return v.Last()
	}

	return new(Chain).Last()
}
