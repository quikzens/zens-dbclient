package httpserver

import (
	"encoding/json"
	"net/http"
	"zens-db/entity"
	"zens-db/helper"
)

func (h *Handler) GetConnections(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	connections := h.usecase.GetConnections(ctx)

	h.writeSuccess(w, serializeArray(connections, connectionResponse), MetaResponse{
		HTTPCode: http.StatusOK,
	})
}

type createConnectionRequest struct {
	Host         string `json:"host" validate:"required"`
	Port         string `json:"port" validate:"required"`
	DatabaseName string `json:"database_name" validate:"required"`
	User         string `json:"user" validate:"required"`
	Password     string `json:"password" validate:"required"`
}

func (h *Handler) CreateConnection(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var request createConnectionRequest
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		h.writeError(w, entity.JSONBadRequestError{})
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

	h.writeSuccess(w, result, MetaResponse{
		HTTPCode: http.StatusOK,
	})
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

	h.writeSuccess(w, result, MetaResponse{
		HTTPCode: http.StatusOK,
	})
}

func connectionResponse(c entity.Connection) interface{} {
	return struct {
		ConnectionID int    `json:"connection_id"`
		Host         string `json:"host"`
		Port         string `json:"port"`
		User         string `json:"user"`
		Password     string `json:"password"`
		DatabaseName string `json:"database_name"`
	}{
		ConnectionID: c.Id,
		Host:         c.Credential.Host,
		Port:         c.Credential.Port,
		User:         c.Credential.User,
		Password:     c.Credential.Password,
		DatabaseName: c.Credential.DatabaseName,
	}
}
