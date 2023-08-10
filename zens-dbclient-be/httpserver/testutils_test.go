package httpserver

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
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
