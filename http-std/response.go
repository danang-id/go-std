package http_std

import (
	"time"
)

// Response defines standard HTTP Response.
type Response struct {
	Success   bool            `json:"success"`
	Timestamp time.Time       `json:"timestamp,omitempty"`
	Message   string          `json:"message,omitempty"`
	Errors    []ResponseError `json:"errors,omitempty"`
	Data      interface{}     `json:"data,omitempty"`
}

// EmptyResponse creates a new empty HTTP Response.
func EmptyResponse() *Response {
	return &Response{
		Success:   true,
		Timestamp: time.Now(),
		Message:   "",
		Errors:    nil,
		Data:      nil,
	}
}

// GetMessage returns the HTTP Response's message and whether it's exist or not.
func (response *Response) GetMessage() (message string, exists bool) {
	return response.Message, len(response.Message) > 0
}

// GetErrors returns the HTTP Response's errors and whether it's exist or not.
func (response *Response) GetErrors() (errors []ResponseError, exists bool) {
	return response.Errors, response.Errors != nil && len(response.Errors) > 0
}

// GetData returns the HTTP Response's data and whether it's exist or not.
func (response *Response) GetData() (data interface{}, exists bool) {
	return response.Data, response.Data != nil
}
