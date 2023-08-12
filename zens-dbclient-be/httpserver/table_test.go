package httpserver

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"testing"
	"zens-db/entity"

	"github.com/stretchr/testify/assert"
)

var (
	validTableName   = "valid_table_name"   // ensure this table is exist (and has columns) in your database
	invalidTableName = "invalid_table_name" // ensure this table is not exist in your database
)

type getTablesHttpResponse struct {
	testResponse
	Data getTablesResponse `json:"data"`
}

func TestGetTables(t *testing.T) {
	ts := newTestServer(t)
	defer ts.Close()

	_, _, createConnectionResp := createConnection(ts, t, createConnectionRequests["valid"])
	connectionId := strconv.FormatInt(int64(createConnectionResp.Data.ConnectionId), 10)

	tests := []struct {
		name         string
		connectionId string
		wantCode     int
	}{
		{
			name:         "Existing Connection",
			connectionId: connectionId,
			wantCode:     http.StatusOK,
		},
		{
			name:         "Non-Existing Connection",
			connectionId: "100",
			wantCode:     http.StatusNotFound,
		},
		{
			name:         "Text Connection ID",
			connectionId: "non-number-string",
			wantCode:     http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			code, _, body := ts.get(t, fmt.Sprintf("/%s/tables", tt.connectionId))

			// check status code
			assert.Equal(t, tt.wantCode, code)

			// parse response
			var resp getTablesHttpResponse
			err := json.Unmarshal([]byte(body), &resp)
			assert.Nil(t, err)
		})
	}
}

type getTableColumnsHttpResponse struct {
	testResponse
	Data getTableColumnsResponse `json:"data"`
}

func TestGetTableColumns(t *testing.T) {
	ts := newTestServer(t)
	defer ts.Close()

	_, _, createConnectionResp := createConnection(ts, t, createConnectionRequests["valid"])
	connectionId := strconv.FormatInt(int64(createConnectionResp.Data.ConnectionId), 10)

	tests := []struct {
		name              string
		connectionId      string
		tableName         string
		wantCode          int
		withRespDataCheck bool
		respDataCheck     func(*testing.T, getTableColumnsResponse)
	}{
		{
			name:              "Existing Connection and Table",
			connectionId:      connectionId,
			tableName:         validTableName,
			wantCode:          http.StatusOK,
			withRespDataCheck: true,
			respDataCheck: func(t *testing.T, data getTableColumnsResponse) {
				assert.NotEqual(t, 0, len(data))
				for _, column := range data {
					assert.NotEqual(t, "", column["column_name"])
					assert.NotEqual(t, "", column["data_type"])
				}
			},
		},
		{
			name:         "Non-Existing Connection",
			connectionId: "100",
			wantCode:     http.StatusNotFound,
		},
		{
			name:         "Text Connection ID",
			connectionId: "non-number-string",
			wantCode:     http.StatusBadRequest,
		},
		{
			name:              "Non-Existing Table",
			connectionId:      connectionId,
			tableName:         invalidTableName,
			wantCode:          http.StatusOK,
			withRespDataCheck: true,
			respDataCheck: func(t *testing.T, data getTableColumnsResponse) {
				assert.Equal(t, 0, len(data))
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			code, _, body := ts.get(t, fmt.Sprintf("/%s/tables/%s/columns", tt.connectionId, tt.tableName))

			// check status code
			assert.Equal(t, tt.wantCode, code)

			// parse response
			var resp getTableColumnsHttpResponse
			err := json.Unmarshal([]byte(body), &resp)
			assert.Nil(t, err)

			// check response
			if tt.withRespDataCheck {
				tt.respDataCheck(t, resp.Data)
			}
		})
	}
}

type getTableRecordsHttpResponse struct {
	testResponse
	Data getTableRecordsResponse `json:"data"`
}

func TestGetTableRecords(t *testing.T) {
	ts := newTestServer(t)
	defer ts.Close()

	_, _, createConnectionResp := createConnection(ts, t, createConnectionRequests["valid"])
	connectionId := strconv.FormatInt(int64(createConnectionResp.Data.ConnectionId), 10)

	tests := []struct {
		name         string
		connectionId string
		tableName    string
		wantCode     int
	}{
		{
			name:         "Existing Connection and Table",
			connectionId: connectionId,
			tableName:    validTableName,
			wantCode:     http.StatusOK,
		},
		{
			name:         "Non-Existing Connection",
			connectionId: "100",
			wantCode:     http.StatusNotFound,
		},
		{
			name:         "Text Connection ID",
			connectionId: "non-number-string",
			wantCode:     http.StatusBadRequest,
		},
		{
			name:         "Non-Existing Table",
			connectionId: connectionId,
			tableName:    invalidTableName,
			wantCode:     http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			code, _, body := ts.postJson(t, fmt.Sprintf("/%s/tables/%s/records", tt.connectionId, tt.tableName), getTableRecordsRequest{
				Conditions: []entity.Condition{},
			})

			// check status code
			assert.Equal(t, tt.wantCode, code)

			// parse response
			var resp getTableRecordsHttpResponse
			err := json.Unmarshal([]byte(body), &resp)
			assert.Nil(t, err)
		})
	}
}
