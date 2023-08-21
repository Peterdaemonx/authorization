package iso8583

import (
	"fmt"
)

func NewElementError(number int, err error) error {
	return ElementError{
		Number: number,
		Err:    err,
	}
}

type ElementError struct {
	Number int
	Err    error
}

func (e ElementError) Error() string {
	//nolint:errorlint
	if subelement, ok := e.Err.(ElementError); ok {
		return fmt.Sprintf("DE %d subfield %d: %s", e.Number, subelement.Number, subelement.Err)
	}

	return fmt.Sprintf("DE %d: %s", e.Number, e.Err)
}
