package entity

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"gitlab.cmpayments.local/creditcard/authorization/internal/processing/cardinfo"
)

const (
	firstVisibleDigits = 8
	lastVisibleDigits  = 4
)

type Card struct {
	Number       string
	Cvv          string
	Holder       string
	Expiry       Expiry
	MaskedPan    string
	PanTokenID   string
	Info           cardinfo.Range
	SequenceNumber string
	Track2Data     string
}

func (c Card) CardholderName() string {
	return c.Holder
}

func (c Card) CardNumber() string {
	return c.Number
}

// IsMultiCountry 59 can be more than one country. So you'll send the country with it.
func (c Card) IsMultiCountry() bool {
	return strings.HasPrefix(c.MaskedPan, "59")
}

type Expiry struct {
	Year  string
	Month string
}

func (e Expiry) MustYearToInt() int {
	year, _ := strconv.Atoi(e.Year)
	century := time.Now().Year() / 100 * 100 // results into the century. example: 20
	return century + year
}

func (e Expiry) MustMonthToInt() int {
	month, _ := strconv.Atoi(e.Month)
	return month
}

func (e Expiry) String() string {
	return fmt.Sprintf("%02s%02s", e.Year, e.Month)
}

func MaskPan(pan string) string {
	if len(pan) <= 10 {
		return strings.Repeat("#", len(pan))
	}

	return pan[:firstVisibleDigits] + strings.Repeat("#", len(pan)-firstVisibleDigits-lastVisibleDigits) + pan[len(pan)-lastVisibleDigits:]
}
