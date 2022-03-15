package httpstd

// ResponseError defines the error information may be contained in the httpstd.Response.
type ResponseError struct {
	Code   int64  `json:"code"`
	Reason string `json:"reason"`
}
