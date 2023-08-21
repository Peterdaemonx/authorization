package data

import (
	"errors"
	"strconv"
	"time"
)

var ErrInvalidLocalTransactionDateTime = errors.New("invalid transaction date time")

type LocalTransactionDateTime time.Time

func (lt *LocalTransactionDateTime) UnmarshalJSON(b []byte) error {
	unquotedJSONValue, err := strconv.Unquote(string(b))
	if err != nil {
		return ErrInvalidLocalTransactionDateTime
	}

	t, err := time.Parse("2006-01-02 15:04:05", unquotedJSONValue)
	if err != nil {
		return err
	}

	*lt = LocalTransactionDateTime(t)

	return nil
}

func (lt *LocalTransactionDateTime) MarshalJSON() ([]byte, error) {
	jsonValue := lt.Format("2006-01-02 15:04:05")

	quotedJSONValue := strconv.Quote(jsonValue)

	return []byte(quotedJSONValue), nil
}

func (lt *LocalTransactionDateTime) Format(s string) string {
	t := time.Time(*lt)
	return t.Format(s)
}
