package web

import (
	"net/url"
	"strconv"
	"time"

	"gitlab.cmpayments.local/creditcard/platform/http/validator"
)

func ReadInt(qs url.Values, key string, defaultValue int, v *validator.Validator) int {
	s := qs.Get(key)

	if s == "" {
		return defaultValue
	}

	i, err := strconv.Atoi(s)
	if err != nil {
		v.AddError(key, []string{"must be an integer value"})
		return defaultValue
	}

	return i
}

func ReadString(qs url.Values, key string, defaultValue string) string {
	s := qs.Get(key)

	if s == "" {
		return defaultValue
	}

	return s
}

func ReadDate(qs url.Values, key string, defaultValue time.Time, v *validator.Validator) time.Time {
	s := qs.Get(key)

	if s == "" {
		return defaultValue
	}

	t, err := time.Parse("2006-01-02", s)
	if err != nil {
		v.AddError(key, []string{"date must be in format of yyyy-mm-dd"})
		return defaultValue
	}

	return t
}
