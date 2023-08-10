package httpserver

import "net/http"

func (h *Handler) Health(w http.ResponseWriter, r *http.Request) {
	_, _ = w.Write([]byte("ok"))
}
