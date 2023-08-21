package base1

type F044_AdditionalResponseData struct {
	_                    string `iso8583:"1=ans-1, dataenc=ebcdic, omitempty, minlength=0"`
	_                    string `iso8583:"2=an-1, dataenc=ebcdic, omitempty, minlength=0"`
	_                    string `iso8583:"3=b-1, omitempty, minlength=0"`
	_                    string `iso8583:"4=b-1, omitempty, minlength=0"`
	_                    string `iso8583:"5=ans-1, dataenc=ebcdic, omitempty, minlength=0"`
	_                    string `iso8583:"6=ans-2, dataenc=ebcdic, omitempty, minlength=0"`
	_                    string `iso8583:"7=ans-1, dataenc=ebcdic, omitempty, minlength=0"`
	_                    string `iso8583:"8=ans-1, dataenc=ebcdic, omitempty, minlength=0"`
	_                    string `iso8583:"9=b-1, omitempty, minlength=0"`
	SF10_CVV2ResultCode  string `iso8583:"10=ans-1"`
	_                    string `iso8583:"11=b-2, omitempty"`
	_                    string `iso8583:"12=b-1, omitempty"`
	SF13_CavvResultsCode string `iso8583:"13=b-1, omitempty"`
}
