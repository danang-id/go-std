package http_std

import (
	"encoding/json"
	"net/http"
)

type SendFunc func(statusCode int, value interface{}) error

// ResponseBuilder defines methods to build HTTP Response.
type ResponseBuilder interface {
	// AppendError add an error to HTTP Response.
	AppendError(code int64, reason string) ResponseBuilder

	// AppendErrors add multiple errors to HTTP Response.
	AppendErrors(errors ...ResponseError) ResponseBuilder

	// Build returns the HTTP Response.
	Build() Response

	// Clean all the information set to HTTP Response.
	Clean() ResponseBuilder

	// SetData set the HTTP Response's Data.
	SetData(data interface{}) ResponseBuilder

	// SetErrors set the HTTP Response's Errors.
	SetErrors(errors ...ResponseError) ResponseBuilder

	// SetMessage set the HTTP Response's Message.
	SetMessage(message string) ResponseBuilder

	// SetStatusCode set the status code of HTTP Response.
	SetStatusCode(statusCode int) ResponseBuilder

	// Send the HTTP Response via the SendFunc callback provided.
	Send(callback SendFunc) error

	// Write the HTTP Response to the http-std.ResponseWriter.
	Write(writer http.ResponseWriter) (int, error)
}

type responseBuilderContext struct {
	response   *Response
	statusCode int
}

// NewResponse create a new ResponseBuilder instance.
func NewResponse() ResponseBuilder {
	return &responseBuilderContext{
		response:   EmptyResponse(),
		statusCode: 200,
	}
}

// AppendError add an error to HTTP Response.
func (context *responseBuilderContext) AppendError(code int64, reason string) ResponseBuilder {
	return context.AppendErrors(ResponseError{Code: code, Reason: reason})
}

// AppendErrors add multiple errors to HTTP Response.
func (context *responseBuilderContext) AppendErrors(errors ...ResponseError) ResponseBuilder {
	if context.response == nil {
		context.Clean()
	}

	if context.response.Errors == nil {
		return context.SetErrors(errors...)
	} else {
		return context.SetErrors(append(context.response.Errors, errors...)...)
	}
}

// Build returns the HTTP Response.
func (context *responseBuilderContext) Build() (response Response) {
	response = *context.response
	context.response = nil
	return response
}

// Clean all the information set to HTTP Response.
func (context *responseBuilderContext) Clean() ResponseBuilder {
	context.response = EmptyResponse()
	return context
}

// SetData set the HTTP Response's Data.
func (context *responseBuilderContext) SetData(data interface{}) ResponseBuilder {
	if context.response == nil {
		context.Clean()
	}

	context.response.Success = true
	context.response.Data = data
	return context
}

// SetErrors set the HTTP Response's Errors.
func (context *responseBuilderContext) SetErrors(errors ...ResponseError) ResponseBuilder {
	if context.response == nil {
		context.Clean()
	}

	context.response.Success = false
	context.response.Errors = errors
	return context
}

// SetMessage set the HTTP Response Response's Message.
func (context *responseBuilderContext) SetMessage(message string) ResponseBuilder {
	if context.response == nil {
		context.Clean()
	}

	context.response.Message = message
	return context
}

// SetStatusCode set the status code of HTTP Response.
func (context *responseBuilderContext) SetStatusCode(statusCode int) ResponseBuilder {
	context.statusCode = statusCode
	return context
}

// Send the HTTP Response via the SendFunc callback provided.
func (context *responseBuilderContext) Send(callback SendFunc) error {
	return callback(context.statusCode, context.response)
}

// Write the HTTP Response to the http-std.ResponseWriter.
func (context *responseBuilderContext) Write(writer http.ResponseWriter) (int, error) {
	writer.WriteHeader(context.statusCode)
	jsonBody, err := json.Marshal(context.response)
	if err != nil {
		return 0, err
	}

	return writer.Write(jsonBody)
}
