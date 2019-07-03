package errors

import (
	"errors"
)

var (
	TimeStampError = GenerateError("time should be a unix timestamp")
	UnauthorizedError = GenerateError("you are not authorized yo perform this action")

)

func GenerateError(err string) error {
	return errors.New(err)
}
