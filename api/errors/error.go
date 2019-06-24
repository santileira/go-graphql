package errors

import (
	"errors"
)

var (
	TimeStampError = GenerateError("time should be a unix timestamp")
)

func GenerateError(err string) error {
	return errors.New(err)
}
