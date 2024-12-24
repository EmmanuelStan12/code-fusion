package errors

import "net/http"

type CustomError struct {
	Type    string
	Message string
	Code    int
	Err     error
	Params  []string
}

func NewError(errType string, err error, code int) *CustomError {
	return &CustomError{
		Type: errType,
		Err:  err,
		Code: code,
	}
}

func BadRequest(errType string, err error) *CustomError {
	return &CustomError{
		Type: errType,
		Err:  err,
		Code: http.StatusBadRequest,
	}
}

func InternalServerError(errType string, err error) *CustomError {
	return &CustomError{
		Type: errType,
		Err:  err,
		Code: http.StatusInternalServerError,
	}
}

func Unauthorized(errType string, err error) *CustomError {
	return &CustomError{
		Type: errType,
		Err:  err,
		Code: http.StatusUnauthorized,
	}
}

func ValidationError(errType string, keyValues ...string) *CustomError {
	return &CustomError{
		Type:   errType,
		Params: keyValues,
		Code:   http.StatusBadRequest,
	}
}
