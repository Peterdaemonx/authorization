package entity

import (
	"fmt"
)

type Recurring struct {
	Initial    bool
	Subsequent bool
	TraceID    string
}

type MTraceID struct {
	FinancialNetworkCode   string
	BanknetReferenceNumber string
	NetworkReportingDate   string
}

func (t MTraceID) String() string {
	return fmt.Sprintf("%s%s%s", t.FinancialNetworkCode, t.BanknetReferenceNumber, t.NetworkReportingDate)
}

func (t MTraceID) NetworkData() string {
	return fmt.Sprintf("%s%s", t.FinancialNetworkCode, t.BanknetReferenceNumber)
}

func TraceIDFromString(traceID string) MTraceID {
	if traceID != "" {
		return MTraceID{
			FinancialNetworkCode:   traceID[:3],
			BanknetReferenceNumber: traceID[3:9],
			NetworkReportingDate:   traceID[9:],
		}
	}
	return MTraceID{}
}
