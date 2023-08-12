package httpserver

import (
	"encoding/json"
	"net/http"
	"zens-db/entity"
	"zens-db/helper"

	"github.com/go-chi/chi/v5"
)

type getTablesResponse []string

func (h *Handler) GetTables(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	connectionId, err := helper.GetUrlIntParam(r, "connection_id", "connection_id must be an integer")
	if err != nil {
		h.writeError(w, err)
		return
	}

	tableNames, err := h.usecase.GetTables(ctx, connectionId)
	if err != nil {
		h.writeError(w, err)
		return
	}

	resp := tableNames
	h.writeSuccess(w, resp, MetaResponse{
		HTTPCode: http.StatusOK,
	})
}

type getTableColumnsResponse []map[string]any

func (h *Handler) GetTableColumns(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	tableName := chi.URLParam(r, "table_name")
	connectionId, err := helper.GetUrlIntParam(r, "connection_id", "connection_id must be an integer")
	if err != nil {
		h.writeError(w, err)
		return
	}

	tableColumns, err := h.usecase.GetTableColumns(ctx, connectionId, tableName)
	if err != nil {
		h.writeError(w, err)
		return
	}

	resp := tableColumns
	h.writeSuccess(w, resp, MetaResponse{
		HTTPCode: http.StatusOK,
	})
}

type getTableRecordsRequest struct {
	Conditions []entity.Condition `json:"conditions"`
}

type getTableRecordsResponse []map[string]any

func (h *Handler) GetTableRecords(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	queryParams := r.URL.Query()
	tableName := chi.URLParam(r, "table_name")
	connectionId, err := helper.GetUrlIntParam(r, "connection_id", "connection_id must be an integer")
	if err != nil {
		h.writeError(w, err)
		return
	}

	var request getTableRecordsRequest
	err = json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		h.writeError(w, entity.JSONBadRequestError{})
		return
	}

	limit, err := helper.GetQueryIntParam(queryParams, "limit", "limit must be a number")
	if err != nil {
		h.writeError(w, err)
		return
	}

	offset, err := helper.GetQueryIntParam(queryParams, "offset", "offset must be a number")
	if err != nil {
		h.writeError(w, err)
		return
	}

	result, err := h.usecase.GetTableRecords(ctx, connectionId, entity.GetTableRecordsParam{
		TableName:  tableName,
		Limit:      limit,
		Offset:     offset,
		SortBy:     queryParams.Get("sort_by"),
		OrderBy:    queryParams.Get("order_by"),
		Conditions: request.Conditions,
	})
	if err != nil {
		h.writeError(w, err)
		return
	}

	resp := result.Data
	h.writeSuccess(w, resp, MetaResponse{
		HTTPCode: http.StatusOK,
		Total:    result.Total,
		Limit:    result.Limit,
		Offset:   result.Offset,
	})
}
