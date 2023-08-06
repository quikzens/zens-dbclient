package entity

type HttpError interface {
	ToHttpError() (int, HttpResponseError)
}

type HttpResponseError struct {
	Field   string `json:"field,omitempty"`
	Message string `json:"message"`
}
