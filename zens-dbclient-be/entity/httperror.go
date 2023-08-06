package entity

import "net/http"

type BadRequestError struct {
	Field   string
	Message string
}

func (e BadRequestError) Error() string {
	return e.Message
}

func (e BadRequestError) ToHttpError() (int, HttpResponseError) {
	return http.StatusBadRequest, HttpResponseError{
		Field:   e.Field,
		Message: e.Message,
	}
}

type JSONBadRequestError struct{}

func (e JSONBadRequestError) Error() string {
	return "JSON is not valid"
}

func (e JSONBadRequestError) ToHttpError() (int, HttpResponseError) {
	return http.StatusBadRequest, HttpResponseError{
		Message: e.Error(),
	}
}

type InternalServerError struct {
	Message string
}

func (e InternalServerError) Error() string {
	return e.Message
}

func (e InternalServerError) ToHttpError() (int, HttpResponseError) {
	return http.StatusInternalServerError, HttpResponseError{
		Message: e.Message,
	}
}
