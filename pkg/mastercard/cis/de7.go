package cis

import (
	"fmt"
	"time"

	"gitlab.cmpayments.local/creditcard/authorization/pkg/mastercard"
)

// DE 7 (Transmission Date and Time) is the  month, day, and time at which the
// transaction was sent to the Mastercard Network
type DE7_TransmissionDateAndTime struct {
	SF1_Date string `iso8583:"1=n-4"` // MMDD
	SF2_Time string `iso8583:"2=n-6"` // hhmmss
}

func (de *DE7_TransmissionDateAndTime) Time() (time.Time, error) {
	return mastercard.TimeParse(`YYYYMMDDhhmmss`,
		fmt.Sprintf("%d%s%s", time.Now().Year(), de.SF1_Date, de.SF2_Time)) //nolint:wrapcheck
}

func DE7FromTime(t time.Time) *DE7_TransmissionDateAndTime {
	return &DE7_TransmissionDateAndTime{
		SF1_Date: mastercard.TimeFormat(t, `MMDD`),
		SF2_Time: mastercard.TimeFormat(t, `hhmmss`),
	}
}

func DE7FromString(s string) *DE7_TransmissionDateAndTime {
	return &DE7_TransmissionDateAndTime{
		SF1_Date: s[0:4],
		SF2_Time: s[4:],
	}
}

func (de *DE7_TransmissionDateAndTime) String() string {
	return de.SF1_Date + de.SF2_Time
}
