package errors

import (
	"encoding/json"
	"fmt"
	"net/http"
)

var (
	ErrInternalSystem = New(http.StatusInternalServerError, "", "00100", "internal system error")
	ErrBadRequest     = New(http.StatusBadRequest, "", "00101", "bad request")
	ErrValidation     = New(http.StatusBadRequest, "", "00102", "validation error")
	ErrNotFound       = New(http.StatusNotFound, "", "00103", "not found")
	ErrUnauthorized   = New(http.StatusUnauthorized, "", "00104", "unauthorized")
	ErrForbidden      = New(http.StatusForbidden, "", "00105", "access forbidden")
)

type ErrorFields map[string]string

type AppError struct {
	Err           error       `json:"-"`
	Message       string      `json:"message,omitempty"`
	Code          string      `json:"code,omitempty"`
	TransportCode int         `json:"-"`
	Fields        ErrorFields `json:"fields,omitempty"`
}

func (e *AppError) WithFields(fields ErrorFields) {
	e.Fields = fields
}

func (e *AppError) Error() string {
	err := e.Err.Error()

	if len(e.Fields) > 0 {
		for k, v := range e.Fields {
			err += ", " + k + " " + v
		}
	}
	return err
}

func New(transportCode int, service, code, message string) *AppError {
	return &AppError{
		Err:           fmt.Errorf(message),
		Code:          service + "-" + code,
		TransportCode: transportCode,
		Message:       message,
	}
}

func Wrap(err error, service, code, message string) *AppError {
	return &AppError{
		Err:     err,
		Code:    service + "-" + code,
		Message: message,
	}
}

func (e *AppError) Unwrap() error { return e.Err }

func (e *AppError) Marshal() []byte {
	marshal, err := json.Marshal(e)
	if err != nil {
		return nil
	}
	return marshal
}
