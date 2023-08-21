package data

import (
	"database/sql"
	"errors"
	"fmt"
	"strconv"
)

var ErrInvalidEcommerceIndicator = errors.New("invalid ecommerce indicator")

//nolint:deadcode,varcheck,unused
const (
	mcauth3dsfail      EcommerceIndicator = 0
	mcbankorcardnot3ds EcommerceIndicator = 1
	mcauth3dsok        EcommerceIndicator = 2
	auth3dsok          EcommerceIndicator = 5
	bankorcardnot3ds   EcommerceIndicator = 6
	auth3dsfail        EcommerceIndicator = 7
)

type EcommerceIndicator int

func (e *EcommerceIndicator) UnmarshalJSON(b []byte) error {
	unquotedJSONValue, err := strconv.Unquote(string(b))
	if err != nil {
		return ErrInvalidEcommerceIndicator
	}

	i, err := strconv.ParseInt(unquotedJSONValue, 10, 32)
	if err != nil {
		return ErrInvalidEcommerceIndicator
	}

	*e = EcommerceIndicator(i)

	return nil
}

func (e *EcommerceIndicator) MarshalJSON() ([]byte, error) {
	jsonValue := fmt.Sprintf("%0d", *e)

	quotedJSONValue := strconv.Quote(jsonValue)

	return []byte(quotedJSONValue), nil
}

func (e *EcommerceIndicator) String() string {
	return fmt.Sprintf("%d", *e)
}

func (e *EcommerceIndicator) NullInt64() sql.NullInt64 {
	if *e == 0 {
		return sql.NullInt64{}
	}

	return sql.NullInt64{
		Int64: int64(*e),
		Valid: true,
	}
}
