package base1

import (
	"fmt"
	"time"

	"gitlab.cmpayments.local/creditcard/authorization/pkg/mastercard"
)

// DE 7 (Transmission Date and Time) is the  month, day, and time at which the
// transaction was sent to the Mastercard Network
type F007_TransmissionDateAndTime struct {
	SF1_Date string `iso8583:"1=n-4"` // MMDD
	SF2_Time string `iso8583:"2=n-6"` // hhmmss
}

func (de *F007_TransmissionDateAndTime) Time() (time.Time, error) {
	return mastercard.TimeParse(`YYYYMMDDhhmmss`,
		fmt.Sprintf("%d%s%s", time.Now().Year(), de.SF1_Date, de.SF2_Time)) //nolint:wrapcheck
}

func F007FromTime(t time.Time) F007_TransmissionDateAndTime {
	return F007_TransmissionDateAndTime{
		SF1_Date: mastercard.TimeFormat(t, `MMDD`),
		SF2_Time: mastercard.TimeFormat(t, `hhmmss`),
	}
}

func F007FromString(s string) *F007_TransmissionDateAndTime {
	return &F007_TransmissionDateAndTime{
		SF1_Date: s[0:4],
		SF2_Time: s[4:],
	}
}

func (de *F007_TransmissionDateAndTime) String() string {
	return de.SF1_Date + de.SF2_Time
}
