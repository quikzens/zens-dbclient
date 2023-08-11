package httpserver

import (
	"encoding/json"
	"net/http"

	"zens-db/entity"
	"zens-db/usecase"
)

type Handler struct {
	usecase *usecase.Usecase
}

func NewHandler(usecase *usecase.Usecase) *Handler {
	return &Handler{
		usecase: usecase,
	}
}

type MetaResponse struct {
	HTTPCode int      `json:"http_code"`
	Limit    int      `json:"limit,omitempty"`
	Offset   int      `json:"offset,omitempty"`
	Total    int      `json:"total,omitempty"`
	Menus    []string `json:"menus,omitempty"`
}

type Response struct {
	Message          string                    `json:"message,omitempty"`
	Data             interface{}               `json:"data,omitempty"`
	Error            *entity.HttpResponseError `json:"error,omitempty"`
	ValidationErrors []entity.ValidationError  `json:"validation_errors,omitempty"`
	Meta             MetaResponse              `json:"meta"`
}

func (h *Handler) writeSuccessWithMessage(w http.ResponseWriter, data interface{}, message string, meta MetaResponse) {
	res := Response{
		Data:    data,
		Message: message,
		Meta:    meta,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(meta.HTTPCode)
	responseBody, _ := json.Marshal(res)
	_, _ = w.Write(responseBody)
}

func (h *Handler) writeSuccess(w http.ResponseWriter, data interface{}, meta MetaResponse) {
	res := Response{
		Data: data,
		Meta: meta,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(meta.HTTPCode)
	responseBody, _ := json.Marshal(res)
	_, _ = w.Write(responseBody)
}

func (h *Handler) writeError(w http.ResponseWriter, err error) {
	statusCode, httpError := h.translateError(err)
	meta := MetaResponse{
		HTTPCode: statusCode,
	}
	res := Response{
		Error: &httpError,
		Meta:  meta,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(meta.HTTPCode)
	responseBody, _ := json.Marshal(res)
	_, _ = w.Write(responseBody)
}

func (h *Handler) writeValidationError(w http.ResponseWriter, errs map[string]string) {
	validationErrs := make([]entity.ValidationError, 0)
	for field, message := range errs {
		validationErrs = append(validationErrs, entity.ValidationError{
			Field:   field,
			Message: message,
		})
	}
	meta := MetaResponse{
		HTTPCode: http.StatusUnprocessableEntity,
	}
	res := Response{
		ValidationErrors: validationErrs,
		Meta:             meta,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(meta.HTTPCode)
	responseBody, _ := json.Marshal(res)
	_, _ = w.Write(responseBody)
}

func (h *Handler) translateError(err error) (int, entity.HttpResponseError) {
	switch origErr := err.(type) {
	case entity.HttpError:
		return origErr.ToHttpError()
	default:
		return entity.InternalServerError{Message: err.Error()}.ToHttpError()
	}
}

func serializeArray[Input any, Output any](array []Input, serializeFunc func(source Input) Output) []Output {
	res := make([]Output, 0, len(array))
	for _, a := range array {
		res = append(res, serializeFunc(a))
	}

	return res
}
