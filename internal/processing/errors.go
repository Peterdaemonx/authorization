package processing

import (
	"gitlab.cmpayments.local/creditcard/platform/http/errors"
)

const (
	PspNotFound = iota + 100008000
	TransactionNotFoundErrorCode
	InvalidTransactionStateErrorCode
	DetokenizationErrorCode
	MappingRequestErrorCode
	MappingResponseErrorCode
	SchemeCommunicationErrorCode
	PersistErrorCode
)

const (
	// FormatError 30 (Format Error) in the Authorization Request Response/0110 message
	// page 352 CustomerInterfaceSpecification.PDF
	FormatError = "30"
)

var (
	ErrAuthorizationNotFound = errors.ProcessorError{
		Code:    TransactionNotFoundErrorCode,
		Message: "transaction not found",
		Details: []string{"the specified authorization does not exist for the specified authorization ID"},
	}
	ErrPspNotFound = errors.ProcessorError{
		Code:    PspNotFound,
		Message: "psp not found",
		Details: []string{"the specified psp does not exist"},
	}
	ErrDetokenization = errors.ProcessorError{
		Code:    DetokenizationErrorCode,
		Message: "detokenization process failed",
		Details: []string{"failed to detokenize credit card"},
	}
)
