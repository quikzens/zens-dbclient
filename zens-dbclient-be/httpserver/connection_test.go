package httpserver

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
)

var createConnectionRequests = map[string]createConnectionRequest{
	"valid": {
		// fill this with a valid database credential
		Host:         "localhost",
		Port:         "5432",
		DatabaseName: "valid_database",
		User:         "valid_user",
		Password:     "valid_password",
	},
	"invalid": {
		// fill this with an invalid database credential
		Host:         "localhost",
		Port:         "5432",
		DatabaseName: "invalid_database",
		User:         "invalid_user",
		Password:     "invalid_password",
	},
}

type createConnectionHttpResponse struct {
	testResponse
	Data createConnectionResponse `json:"data"`
}

func createConnection(ts *testServer, t *testing.T, req createConnectionRequest) (int, http.Header, createConnectionHttpResponse) {
	code, header, body := ts.postJson(t, "/connections", req)

	// parse response
	var resp createConnectionHttpResponse
	err := json.Unmarshal([]byte(body), &resp)
	assert.Nil(t, err)

	return code, header, resp
}

func TestCreateConnection(t *testing.T) {
	ts := newTestServer(t)
	defer ts.Close()

	tests := []struct {
		name              string
		reqBody           createConnectionRequest
		wantCode          int
		withRespDataCheck bool
		respDataCheck     func(*testing.T, createConnectionResponse)
	}{
		{
			name:              "Valid Connection",
			reqBody:           createConnectionRequests["valid"],
			wantCode:          http.StatusOK,
			withRespDataCheck: true,
			respDataCheck: func(t *testing.T, data createConnectionResponse) {
				assert.NotEqual(t, 0, data.ConnectionId)
			},
		},
		{
			name:              "Invalid Connection",
			reqBody:           createConnectionRequests["invalid"],
			wantCode:          http.StatusInternalServerError,
			withRespDataCheck: false,
		},
		{
			name:              "Empty Connection Credential",
			reqBody:           createConnectionRequest{},
			wantCode:          http.StatusUnprocessableEntity,
			withRespDataCheck: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			code, _, resp := createConnection(ts, t, tt.reqBody)

			// check status code
			assert.Equal(t, tt.wantCode, code)

			// check response
			if tt.withRespDataCheck {
				tt.respDataCheck(t, resp.Data)
			}
		})
	}
}

type getConnectionsHttpResponse struct {
	testResponse
	Data getConnectionsResponse `json:"data"`
}

func TestGetConnections(t *testing.T) {
	ts := newTestServer(t)
	defer ts.Close()

	// create connections
	connectionCreated := 3
	for i := 0; i < connectionCreated; i++ {
		createConnection(ts, t, createConnectionRequests["valid"])
	}

	code, _, body := ts.get(t, "/connections")

	// check status code
	assert.Equal(t, http.StatusOK, code)

	// parse response
	var resp getConnectionsHttpResponse
	err := json.Unmarshal([]byte(body), &resp)
	assert.Nil(t, err)

	// check response
	assert.Len(t, resp.Data, 3)
}

type deleteConnectionHttpResponse struct {
	testResponse
	Data deleteConnectionResponse `json:"data"`
}

func TestDeleteConnection(t *testing.T) {
	ts := newTestServer(t)
	defer ts.Close()

	// create connection
	_, _, createConnectionResp := createConnection(ts, t, createConnectionRequests["valid"])
	connectionId := strconv.FormatInt(int64(createConnectionResp.Data.ConnectionId), 10)

	tests := []struct {
		name              string
		connectionId      string
		wantCode          int
		withRespDataCheck bool
		respDataCheck     func(*testing.T, deleteConnectionResponse)
	}{
		{
			name:              "Existing Connection",
			connectionId:      connectionId,
			wantCode:          http.StatusOK,
			withRespDataCheck: true,
			respDataCheck: func(t *testing.T, data deleteConnectionResponse) {
				assert.NotEqual(t, 0, data.ConnectionId)
			},
		},
		{
			name:              "Non-Existing Connection",
			connectionId:      connectionId,
			wantCode:          http.StatusNotFound,
			withRespDataCheck: false,
		},
		{
			name:              "Text Connection ID",
			connectionId:      "non-number-string",
			wantCode:          http.StatusBadRequest,
			withRespDataCheck: false,
		},
		{
			name:              "Negative Connection ID",
			connectionId:      "-1",
			wantCode:          http.StatusNotFound,
			withRespDataCheck: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			code, _, body := ts.delete(t, fmt.Sprintf("/connections/%s", tt.connectionId))

			// check status code
			assert.Equal(t, tt.wantCode, code)

			// parse response
			var resp deleteConnectionHttpResponse
			err := json.Unmarshal([]byte(body), &resp)
			assert.Nil(t, err)

			// check response
			if tt.withRespDataCheck {
				tt.respDataCheck(t, resp.Data)
			}
		})
	}
}
