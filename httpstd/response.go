package httpstd

import (
	"time"
)

// Response defines standard httpstd.Response.
type Response struct {
	Success   bool            `json:"success"`
	Timestamp time.Time       `json:"timestamp,omitempty"`
	Message   string          `json:"message,omitempty"`
	Errors    []ResponseError `json:"errors,omitempty"`
	Data      interface{}     `json:"data,omitempty"`
}

// EmptyResponse creates a new empty httpstd.Response.
func EmptyResponse() *Response {
	return &Response{
		Success:   true,
		Timestamp: time.Now(),
		Message:   "",
		Errors:    nil,
		Data:      nil,
	}
}

// GetMessage returns the httpstd.Response's message and whether it's exist or not.
func (response *Response) GetMessage() (message string, exists bool) {
	return response.Message, len(response.Message) > 0
}

// GetErrors returns the httpstd.Response's errors and whether it's exist or not.
func (response *Response) GetErrors() (errors []ResponseError, exists bool) {
	return response.Errors, response.Errors != nil && len(response.Errors) > 0
}

// GetData returns the httpstd.Response's data and whether it's exist or not.
func (response *Response) GetData() (data interface{}, exists bool) {
	return response.Data, response.Data != nil
}
