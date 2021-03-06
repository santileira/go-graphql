package errors

import (
	"errors"
)

var (
	DateTimeError     = GenerateError("time should be a string date")
	UnauthorizedError = GenerateError("you aren't authorized to perform this action")
)

// GenerateError returns error with that string.
func GenerateError(err string) error {
	return errors.New(err)
}
