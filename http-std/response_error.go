package http_std

// ResponseError defines the error information may be contained in the HTTP Response.
type ResponseError struct {
	Code   int64  `json:"code"`
	Reason string `json:"reason"`
}
