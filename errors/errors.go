package errors

import "fmt"

type FaturaError struct {
	Message  string
	Request  interface{}
	Response interface{}
	Cause    error
}

func (e *FaturaError) Error() string {
	if e.Cause != nil {
		return fmt.Sprintf("%s: %v", e.Message, e.Cause)
	}
	return e.Message
}

func NewApiError(msg string, req interface{}, resp interface{}) *FaturaError {
	return &FaturaError{
		Message:  msg,
		Request:  req,
		Response: resp,
	}
}

func NewBadResponseError(msg string, req interface{}, resp interface{}) *FaturaError {
	return &FaturaError{
		Message:  msg,
		Request:  req,
		Response: resp,
	}
}

func NewInvalidFormatError(msg string, data interface{}) *FaturaError {
	return &FaturaError{
		Message: msg,
		Request: data,
	}
}

func NewInvalidArgumentError(msg string, data interface{}) *FaturaError {
	return &FaturaError{
		Message: msg,
		Request: data,
	}
}
