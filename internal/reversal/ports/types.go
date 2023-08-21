package ports

type ReversalResponse struct {
	ID                 string             `json:"id"`
	LogID              string             `json:"logId"`
	AuthorizationID    string             `json:"authorizationId"`
	CardSchemeResponse CardSchemeResponse `json:"cardSchemeResponse"`
}

type CardSchemeResponse struct {
	Status  string `json:"status"`
	Code    string `json:"code"`
	Message string `json:"message"`
	TraceID string `json:"traceId,omitempty"`
}
