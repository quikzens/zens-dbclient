package entity

import "net/http"

type BadRequestError struct {
	Field   string
	Message string
}

func (e BadRequestError) Error() string {
	return e.Message
}

func (e BadRequestError) ToHttpError() (int, []HttpResponseError) {
	return http.StatusBadRequest, []HttpResponseError{
		{
			Field:   e.Field,
			Message: e.Message,
		},
	}
}

type ConflictError struct {
	Field   string
	Message string
}

func (e ConflictError) Error() string {
	return e.Message
}

func (e ConflictError) ToHttpError() (int, []HttpResponseError) {
	return http.StatusConflict, []HttpResponseError{
		{
			Field:   e.Field,
			Message: e.Message,
		},
	}
}

type JSONBadRequestError struct{}

func (e JSONBadRequestError) Error() string {
	return "JSON is not valid"
}

func (e JSONBadRequestError) ToHttpError() (int, []HttpResponseError) {
	return http.StatusBadRequest, []HttpResponseError{
		{
			Message: e.Error(),
		},
	}
}

type InternalRouteError struct {
	URL        string
	Message    interface{}
	StatusCode int
}

func (e InternalRouteError) Error() string {
	return "Internal Route Error"
}

func (e InternalRouteError) ToHttpError() (int, []HttpResponseError) {
	return http.StatusInternalServerError, []HttpResponseError{
		{
			Message: e.Error(),
		},
	}
}

type InternalRouteNotFoundError struct {
	Field string
}

func (e InternalRouteNotFoundError) Error() string {
	return e.Field + " Not Found"
}

func (e InternalRouteNotFoundError) ToHttpError() (int, []HttpResponseError) {
	return http.StatusNotFound, []HttpResponseError{
		{
			Field:   e.Field,
			Message: e.Error(),
		},
	}
}

type InternalServerError struct {
	Message string
}

func (e InternalServerError) Error() string {
	return e.Message
}

func (e InternalServerError) ToHttpError() (int, []HttpResponseError) {
	return http.StatusInternalServerError, []HttpResponseError{
		{
			Message: e.Message,
		},
	}
}

type UnauthorizedError struct {
	Message string
}

func (e UnauthorizedError) Error() string {
	return e.Message
}

func (e UnauthorizedError) ToHttpError() (int, []HttpResponseError) {
	return http.StatusUnauthorized, []HttpResponseError{
		{
			Message: e.Message,
		},
	}
}

type TooManyRequestError struct {
	Message string
}

func (e TooManyRequestError) Error() string {
	return e.Message
}

func (e TooManyRequestError) ToHttpError() (int, []HttpResponseError) {
	return http.StatusTooManyRequests, []HttpResponseError{
		{
			Message: e.Message,
		},
	}
}

type ForbiddenError struct{}

func (e ForbiddenError) Error() string {
	return "Forbidden"
}

func (e ForbiddenError) ToHttpError() (int, []HttpResponseError) {
	return http.StatusForbidden, []HttpResponseError{
		{
			Message: "Forbidden",
		},
	}
}

type UnprocessableError struct {
	Message string
}

func (e UnprocessableError) Error() string {
	return e.Message
}

func (e UnprocessableError) ToHttpError() (int, []HttpResponseError) {
	return http.StatusUnprocessableEntity, []HttpResponseError{
		{
			Message: e.Message,
		},
	}
}

type InternalRouteForbiddenError struct {
	Field   string
	Message string
}

func (e InternalRouteForbiddenError) Error() string {
	return e.Message
}

func (e InternalRouteForbiddenError) ToHttpError() (int, []HttpResponseError) {
	return http.StatusForbidden, []HttpResponseError{
		{
			Field:   e.Field,
			Message: e.Error(),
		},
	}
}
