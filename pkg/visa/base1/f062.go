package base1

type F062_CustomPaymentService struct {
	SF1_AuthorizationCharacteristicsIndicator string `iso8583:"1=an-1, dataenc=ebcdic, omitempty, minlength=0"`
	SF2_TransactionIdentifier                 int    `iso8583:"2=n-15, dataenc=bcd4, justify=right, omitempty"`
	SF3_ValidationCode                        string `iso8583:"3=an-4, dataenc=ebcdic, omitempty"`
	SF4_MarketSpecificDataIdentifier          string `iso8583:"4=an-1, dataenc=ebcdic, omitempty"`
	SF5_Duration                              string `iso8583:"5=n-2, dataenc=bcd4, justify=right, omitempty"`
	SF6_Reserved                              string `iso8583:"6=an-1, dataenc=ebcdic, omitempty"`
	SF7_PurchaseIdentifier                    string `iso8583:"7=an-26, dataenc=ebcdic, omitempty"`
	_                                         string `iso8583:"8=an-1, omitempty, minlength=0"`
	_                                         string `iso8583:"9=an-1, omitempty, minlength=0"`
	_                                         string `iso8583:"10=an-1, omitempty, minlength=0"`
	_                                         string `iso8583:"11=an-1, omitempty, minlength=0"`
	_                                         string `iso8583:"12=an-1, omitempty, minlength=0"`
	_                                         string `iso8583:"13=an-1, omitempty, minlength=0"`
	_                                         string `iso8583:"14=an-1, omitempty, minlength=0"`
	_                                         string `iso8583:"15=an-1, omitempty, minlength=0"`
	SF16_Reserved                             string `iso8583:"16=an-2, dataenc=ebcdic, omitempty"`
	SF17_MastercardInterchangeCompliance      string `iso8583:"17=an-15, dataenc=ebcdic, omitempty"`
	_                                         string `iso8583:"18=an-15, minlength=0, omitempty"`
	_                                         string `iso8583:"19=an-15, minlength=0, omitempty"`
	SF20_MerchantVerificationValue            string `iso8583:"20=an-2, dataenc=bcd4, justify=right, omitempty"`
	SF21_RiskAssesmentScoreAndReasonCodes     string `iso8583:"21=an-4, dataenc=ebcdic, justify=right, omitempty"`
	SF22_RiskAssesmentConditionCodes          string `iso8583:"22=an-6, dataenc=ebcdic, justify=right, omitempty"`
	SF23_ProductID                            string `iso8583:"23=an-2, dataenc=ebcdic, justify=right, omitempty"`
	SF24_ProgramIdentifier                    string `iso8583:"24=an-6, dataenc=ebcdic, justify=right, omitempty"`
	SF25_SpendQualifiedIndicator              string `iso8583:"25=an-1, dataenc=ebcdic, justify=right, omitempty"`
	SF26_AccountStatus                        string `iso8583:"26=an-1, dataenc=ebcdic, justify=right, omitempty"`
}
