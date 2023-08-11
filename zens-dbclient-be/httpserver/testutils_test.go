package httpserver

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"zens-db/entity"
	"zens-db/repository"
	"zens-db/usecase"
)

type testServer struct {
	*httptest.Server
}

func newTestServer(t *testing.T) *testServer {
	r := repository.New()
	u := usecase.New(r)
	h := NewHandler(u)
	router := NewRouter(h)
	ts := httptest.NewServer(router)
	return &testServer{ts}
}

func (ts *testServer) delete(t *testing.T, urlPath string) (int, http.Header, string) {
	req, err := http.NewRequest("DELETE", ts.URL+urlPath, nil)
	if err != nil {
		t.Fatal(err)
	}

	rs, err := ts.Client().Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer rs.Body.Close()

	body, err := io.ReadAll(rs.Body)
	if err != nil {
		t.Fatal(err)
	}
	bytes.TrimSpace(body)

	return rs.StatusCode, rs.Header, string(body)
}

func (ts *testServer) get(t *testing.T, urlPath string) (int, http.Header, string) {
	rs, err := ts.Client().Get(ts.URL + urlPath)
	if err != nil {
		t.Fatal(err)
	}
	defer rs.Body.Close()

	body, err := io.ReadAll(rs.Body)
	if err != nil {
		t.Fatal(err)
	}
	bytes.TrimSpace(body)

	return rs.StatusCode, rs.Header, string(body)
}

func (ts *testServer) postJson(t *testing.T, urlPath string, req interface{}) (int, http.Header, string) {
	jsonReq, _ := json.Marshal(req)
	rs, err := ts.Client().Post(ts.URL+urlPath, "application/json", bytes.NewBuffer(jsonReq))
	if err != nil {
		t.Fatal(err)
	}
	defer rs.Body.Close()

	body, err := io.ReadAll(rs.Body)
	if err != nil {
		t.Fatal(err)
	}
	bytes.TrimSpace(body)

	return rs.StatusCode, rs.Header, string(body)
}

type testResponse struct {
	Message          string                    `json:"message,omitempty"`
	Error            *entity.HttpResponseError `json:"error,omitempty"`
	ValidationErrors []entity.ValidationError  `json:"validation_errors,omitempty"`
	Meta             MetaResponse              `json:"meta"`
}
