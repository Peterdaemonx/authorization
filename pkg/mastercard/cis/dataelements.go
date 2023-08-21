package cis

// All the data elements for CIS that are needed
// Fields that have a 'subfields' comment are complex DE's that currently do not have a specific struct to implement their SF's
//
//nolint:lll
type DataElements struct {
	DE2_PrimaryAccountNumber               string                            `iso8583:"2=n..19, minlength=5"`
	DE3_ProcessingCode                     string                            `iso8583:"3=n-6"` // Subfields - 2 and 3 are always 00 for us
	DE4_TransactionAmount                  int64                             `iso8583:"4=n-12, justify=right"`
	DE6_CardholderBillingAmount            int64                             `iso8583:"6=n-12, justify=right"`
	DE7_TransmissionDateTime               *DE7_TransmissionDateAndTime      `iso8583:"7=n-10"`
	DE10_ConversionRateCardholderBilling   string                            `iso8583:"10=n-8"` // Subfields
	DE11_SystemTraceAuditNumber            string                            `iso8583:"11=n-6"`
	DE12_LocalTransactionTime              string                            `iso8583:"12=n-6"`
	DE13_LocalTransactionDate              string                            `iso8583:"13=n-4"`
	DE14_ExpirationDate                    string                            `iso8583:"14=n-4,omitempty"`
	DE15_SettlementDate                    string                            `iso8583:"15=n-4"`
	DE16_ConversionDate                    string                            `iso8583:"16=n-4"`
	DE18_MerchantType                      string                            `iso8583:"18=n-4"`
	DE20_PrimaryAccountNumberCountryCode   string                            `iso8583:"20=n-3, justify=right"`
	DE22_PointOfServiceEntryMode           string                            `iso8583:"22=n-3"`
	DE28_TransactionFeeAmount              *DE28_TransactionFeeAmount        `iso8583:"28=an-9"`
	DE32_AcquringInstitutionCode           string                            `iso8583:"32=n..6"`
	DE33_ForwardingInstitutionIDCode       string                            `iso8583:"33=n..6"`
	DE35_TrackTwoData                      string                            `iso8583:"35=ans..37"`
	DE37_RetrievalReferenceNumber          string                            `iso8583:"37=an-12"` // Subfields
	DE38_AuthorizationIdResponse           string                            `iso8583:"38=ans-6, justify=left"`
	DE39_ResponseCode                      string                            `iso8583:"39=an-2"`
	DE41_CardAcceptorTerminalId            string                            `iso8583:"41=ans-8"`
	DE42_CardAcceptorCodeId                string                            `iso8583:"42=ans-15, justify=left"`
	DE43_CardAcceptorNameAndLocation       *DE43_CardAcceptorNameAndLocation `iso8583:"43=ans-40"`
	DE44_AdditionalResponseData            string                            `iso8583:"44=ans..25"`
	DE48_AdditionalData                    *DE48_AdditionalData              `iso8583:"48=ans...999"`
	DE49_TransactionCurrencyCode           string                            `iso8583:"49=n-3"`
	DE51_CardholderBillingCurrencyCode     string                            `iso8583:"51=n-3"`
	DE52_PinData                           string                            `iso8583:"52=b-8"`
	DE53_SecurityRelatedControlInformation string                            `iso8583:"53=n-16"` // Subfields
	DE54_AdditionalAmounts                 []DE54_AmountsAdditional          `iso8583:"54=ans...120"`
	DE56_PaymentAccountData                string                            `iso8583:"56=an...37"`
	DE60_AdviceReasonCode                  string                            `iso8583:"60=ans...60"` // Subfields
	DE61_PointOfServiceData                *DE61_PointOfServiceData          `iso8583:"61=ans...26"`
	DE63_NetworkData                       string                            `iso8583:"63=an...50"` // Subfields
	DE70_NetworkManagementInformationCode  string                            `iso8583:"70=n-3"`
	DE90_OriginalDataElements              *DE90_OriginalDataElements        `iso8583:"90=n-42"`
	DE94_ServiceIndicator                  string                            `iso8583:"94=ans-7"`
	DE95_ReplacementAmounts                *DE95_ReplacementAmounts          `iso8583:"95=n-42"`
	DE96_MessageSecurityCode               string                            `iso8583:"96=n-8"`
	DE108_MoneySendReferenceData           string                            `iso8583:"108=ans...999"`
	DE112_AdditionalDataNationalUse        string                            `iso8583:"112=ans...591"`
	DE120_RecordData                       string                            `iso8583:"120=ans...999"` // Subfields
	DE121_AuthorizingAgentIDCode           string                            `iso8583:"121=n...6"`     // Subfields
	DE124_MemberDefinedData                string                            `iso8583:"124=ans...999"`
	DE123_ReceiptFreeText                  string                            `iso8583:"123=ans...512"`
	DE126_PrivateData                      string                            `iso8583:"126=ans...100"`
	DE127_PrivateData                      string                            `iso8583:"127=ans...100"`
}
