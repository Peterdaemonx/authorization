package entity

type ExemptionType string

const (
	MerchantInitiatedExemption       ExemptionType = `merchantInitiated`
	TransactionRiskAnalysisExemption ExemptionType = `transactionRiskAnalysis`
	RecurringExemption               ExemptionType = `recurring`
	LowValueExemption                ExemptionType = `lowValue`
	SCADelegationExemption           ExemptionType = `scaDelegation`
	SecureCorporateExemption         ExemptionType = `secureCorporate`
)

var (
	exemptionTypeMap = map[string]ExemptionType{
		`merchantInitiated`: MerchantInitiatedExemption,
		`recurring`:         RecurringExemption,
		`lowValue`:          LowValueExemption,
	}
)

func IsValidExemption(exemption string) bool {
	_, exists := exemptionTypeMap[exemption]
	return exists
}
