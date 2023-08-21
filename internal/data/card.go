package data

// swagger:enum CardScheme
type CardScheme string

const (
	Visa       CardScheme = "visa"
	Mastercard CardScheme = "mastercard"
	Unknown    CardScheme = "unknown"
)

// swagger:model CardResponse
type CardResponse struct {
	// Card scheme
	//
	// required: true
	Scheme CardScheme `json:"scheme"`
	// Card number
	//
	// required: true
	// minimum length: 9
	// maximum length: 19
	// pattern: ^\\d+$
	// example: 222300######2704
	Number string `json:"number"`
}
