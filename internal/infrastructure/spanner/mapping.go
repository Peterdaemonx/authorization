package spanner

import (
	"database/sql"
	"time"
)

func Ptr[T any](v T) *T {
	return &v
}

func NewNullTime(t time.Time) sql.NullTime {
	if t.IsZero() {
		return sql.NullTime{}
	}
	return sql.NullTime{
		Time:  t,
		Valid: true,
	}
}

func NewNullError(e error) sql.NullString {
	if e == nil {
		return sql.NullString{}
	}
	return sql.NullString{
		String: e.Error(),
		Valid:  true,
	}
}
func NewNullString(s string) sql.NullString {
	if len(s) == 0 {
		return sql.NullString{}
	}

	return sql.NullString{
		String: s,
		Valid:  true,
	}
}

func NewNullInt64(i int) sql.NullInt64 {
	if i == 0 {
		return sql.NullInt64{}
	}

	return sql.NullInt64{
		Int64: int64(i),
		Valid: true,
	}
}

func NewNullInt64Pointer(i *int) sql.NullInt64 {
	if i == nil {
		return sql.NullInt64{}
	}

	return sql.NullInt64{
		Int64: int64(*i),
		Valid: true,
	}
}

func MapNullTimeStr(value interface{}) interface{} {
	if value == nil {
		return sql.NullTime{}
	}
	return sql.NullTime{
		Time:  value.(time.Time),
		Valid: true,
	}
}
