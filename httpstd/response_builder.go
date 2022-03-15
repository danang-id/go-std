package httpstd

import (
	"encoding/json"
	"net/http"
)

type JSONEncoder func(value interface{}) ([]byte, error)

// SendFunc defines the function to send the httpstd.Response.
type SendFunc func(statusCode int, value interface{}) error

// ResponseBuilder defines methods to build httpstd.Response.
type ResponseBuilder interface {
	// AppendError add a ResponseError to httpstd.Response.
	AppendError(code int64, reason string) ResponseBuilder

	// AppendErrors add multiple ResponseError to httpstd.Response.
	AppendErrors(errors ...ResponseError) ResponseBuilder

	// Build returns the copy httpstd.Response. It also calls ResponseBuilder.Clean() to clear the httpstd.Response
	// cache in the current ResponseBuilder context.
	Build() Response

	// Clean all the information set to httpstd.Response in the current ResponseBuilder context.
	Clean(statusCode ...int) ResponseBuilder

	// SetData set the httpstd.Response's Data.
	SetData(data interface{}) ResponseBuilder

	// SetErrors set the httpstd.Response's Errors.
	SetErrors(errors ...ResponseError) ResponseBuilder

	// SetJSONEncoder set the JSONEncoder used in the current ResponseBuilder context. This encoder is used to encode
	// the httpstd.Response to JSON when ResponseBuilder.Write() is called. This function should be called if you
	// do not want to use the default encoding/json to encode JSON.
	SetJSONEncoder(encoder JSONEncoder) ResponseBuilder

	// SetMessage set the httpstd.Response's Message.
	SetMessage(message string) ResponseBuilder

	// SetStatusCode set the status code of httpstd.Response.
	SetStatusCode(statusCode int) ResponseBuilder

	// Send will call ResponseBuilder.Build() and then send the resulted httpstd.Response via the provided
	// SendFunc callback.
	Send(callback SendFunc) error

	// Write will call ResponseBuilder.Build() and then write the resulted httpstd.Response using the provided
	// http.ResponseWriter. It will also serialize httpstd.Response into JSON and set the Content-Type header to
	// application/json.
	Write(writer http.ResponseWriter) (int, error)
}

type responseBuilderContext struct {
	encodeJSON        JSONEncoder
	defaultStatusCode int
	statusCode        int
	response          *Response
}

// NewResponse create a new ResponseBuilder instance.
func NewResponse(statusCode ...int) ResponseBuilder {
	defaultStatusCode := http.StatusOK
	if len(statusCode) > 0 {
		defaultStatusCode = statusCode[0]
	}

	return &responseBuilderContext{
		encodeJSON:        json.Marshal,
		defaultStatusCode: defaultStatusCode,
		statusCode:        defaultStatusCode,
		response:          EmptyResponse(),
	}
}

// AppendError add a ResponseError to httpstd.Response.
func (context *responseBuilderContext) AppendError(code int64, reason string) ResponseBuilder {
	return context.AppendErrors(ResponseError{Code: code, Reason: reason})
}

// AppendErrors add multiple ResponseError to httpstd.Response.
func (context *responseBuilderContext) AppendErrors(errors ...ResponseError) ResponseBuilder {
	if context.response.Errors == nil {
		return context.SetErrors(errors...)
	} else {
		return context.SetErrors(append(context.response.Errors, errors...)...)
	}
}

// Build returns the copy httpstd.Response. It also calls ResponseBuilder.Clean() to clear the httpstd.Response
// cache in the current ResponseBuilder context.
func (context *responseBuilderContext) Build() (response Response) {
	response = *context.response
	context.Clean()

	return response
}

// Clean all the information set to httpstd.Response in the current ResponseBuilder context.
func (context *responseBuilderContext) Clean(statusCode ...int) ResponseBuilder {
	defaultStatusCode := context.defaultStatusCode
	if len(statusCode) > 0 {
		defaultStatusCode = statusCode[0]
	}

	context.statusCode = defaultStatusCode
	context.response = EmptyResponse()
	return context
}

// SetData set the httpstd.Response's Data.
func (context *responseBuilderContext) SetData(data interface{}) ResponseBuilder {
	context.response.Data = data
	return context
}

// SetErrors set the httpstd.Response's Errors.
func (context *responseBuilderContext) SetErrors(errors ...ResponseError) ResponseBuilder {
	if len(errors) > 0 {
		context.response.Success = false
	}
	context.response.Errors = errors
	return context
}

// SetJSONEncoder set the JSONEncoder used in the current ResponseBuilder context. This encoder is used to encode
// the httpstd.Response to JSON when ResponseBuilder.Write() is called.
func (context *responseBuilderContext) SetJSONEncoder(encoder JSONEncoder) ResponseBuilder {
	if encoder != nil {
		context.encodeJSON = encoder
	}

	return context
}

// SetMessage set the httpstd.Response's Message.
func (context *responseBuilderContext) SetMessage(message string) ResponseBuilder {
	context.response.Message = message
	return context
}

// SetStatusCode set the status code of httpstd.Response.
func (context *responseBuilderContext) SetStatusCode(statusCode int) ResponseBuilder {
	context.statusCode = statusCode
	return context
}

// Send will call ResponseBuilder.Build() and then send the resulted httpstd.Response via the provided
// SendFunc callback.
func (context *responseBuilderContext) Send(callback SendFunc) error {
	response := context.Build()
	return callback(context.statusCode, response)
}

// Write will call ResponseBuilder.Build() and then write the resulted httpstd.Response using the provided
// http.ResponseWriter. It will also serialize httpstd.Response into JSON and set the Content-Type header to
// application/json.
func (context *responseBuilderContext) Write(writer http.ResponseWriter) (int, error) {
	response := context.Build()
	data, err := context.encodeJSON(response)
	if err != nil {
		return 0, err
	}

	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(context.statusCode)
	return writer.Write(data)
}
