package httpserver

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
)

func NewRouter(h *Handler) chi.Router {
	r := chi.NewRouter()

	r.Use(cors.Handler(cors.Options{
		// AllowedOrigins:   []string{"https://foo.com"}, // use this to allow specific origin hosts
		AllowOriginFunc:  func(r *http.Request, origin string) bool { return true },
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		AllowCredentials: false,
	}))

	r.Get("/health", h.Health)

	r.Route("/connections", func(r chi.Router) {
		r.Get("/", h.GetConnections)
		r.Post("/", h.CreateConnection)
		r.Delete("/{connection_id}", h.DeleteConnection)
	})

	r.Route("/{connection_id}/tables", func(r chi.Router) {
		r.Get("/", h.GetTables)
		r.Get("/{table_name}/columns", h.GetTableColumns)
		r.Post("/{table_name}/records", h.GetTableRecords)
	})

	return r
}
