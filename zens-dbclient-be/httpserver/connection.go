package httpserver

import (
	"encoding/json"
	"net/http"
	"zens-db/entity"
	"zens-db/helper"
	"zens-db/validator"
)

type connectionResponse struct {
	ConnectionID int    `json:"connection_id"`
	Host         string `json:"host"`
	Port         string `json:"port"`
	User         string `json:"user"`
	Password     string `json:"password"`
	DatabaseName string `json:"database_name"`
}

type getConnectionsResponse []connectionResponse

func (h *Handler) GetConnections(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	connections := h.usecase.GetConnections(ctx)

	resp := serializeArray(connections, serializeConnectionResponse)
	h.writeSuccess(w, resp, MetaResponse{
		HTTPCode: http.StatusOK,
	})
}

type createConnectionRequest struct {
	Host         string `json:"host"`
	Port         string `json:"port"`
	DatabaseName string `json:"database_name"`
	User         string `json:"user"`
	Password     string `json:"password"`
}

func (r *createConnectionRequest) validate() (bool, map[string]string) {
	v := validator.New()
	v.Check(validator.StrRequired(r.Host), "host", "host must be provided")
	v.Check(validator.StrRequired(r.Port), "port", "port must be provided")
	v.Check(validator.StrRequired(r.DatabaseName), "database_name", "database_name must be provided")
	v.Check(validator.StrRequired(r.User), "user", "user must be provided")
	v.Check(validator.StrRequired(r.Password), "password", "password must be provided")
	return v.Valid(), v.Errors
}

type createConnectionResponse struct {
	ConnectionId int `json:"connection_id"`
}

func (h *Handler) CreateConnection(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var request createConnectionRequest
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		h.writeError(w, entity.JSONBadRequestError{})
		return
	}

	if isValid, errs := request.validate(); !isValid {
		h.writeValidationError(w, errs)
		return
	}

	result, err := h.usecase.CreateConnection(ctx, entity.CreateConnectionParam{
		Host:         request.Host,
		Port:         request.Port,
		DatabaseName: request.DatabaseName,
		User:         request.User,
		Password:     request.Password,
	})
	if err != nil {
		h.writeError(w, err)
		return
	}

	resp := createConnectionResponse{
		ConnectionId: result.ConnectionId,
	}
	h.writeSuccessWithMessage(w, resp, "success create connection", MetaResponse{
		HTTPCode: http.StatusOK,
	})
}

type deleteConnectionResponse struct {
	ConnectionId int `json:"connection_id"`
}

func (h *Handler) DeleteConnection(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	connectionId, err := helper.GetUrlIntParam(r, "connection_id", "connection_id must be an integer")
	if err != nil {
		h.writeError(w, err)
		return
	}

	result, err := h.usecase.DeleteConnection(ctx, connectionId)
	if err != nil {
		h.writeError(w, err)
		return
	}

	resp := deleteConnectionResponse{
		ConnectionId: result.ConnectionId,
	}
	h.writeSuccessWithMessage(w, resp, "success delete connection", MetaResponse{
		HTTPCode: http.StatusOK,
	})
}

func serializeConnectionResponse(c entity.Connection) connectionResponse {
	return connectionResponse{
		ConnectionID: c.Id,
		Host:         c.Credential.Host,
		Port:         c.Credential.Port,
		User:         c.Credential.User,
		Password:     c.Credential.Password,
		DatabaseName: c.Credential.DatabaseName,
	}
}
