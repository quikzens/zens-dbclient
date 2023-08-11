package httpserver

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHealth(t *testing.T) {
	ts := newTestServer(t)
	defer ts.Close()

	code, _, body := ts.get(t, "/health")

	assert.Equal(t, http.StatusOK, code)
	assert.Equal(t, "ok", string(body))
}
