package config

import (
	"net/http"

	"zens-db/httpserver"
	"zens-db/repository"
	"zens-db/usecase"
)

func InitHttpHandler(addr string) *http.Server {
	r := repository.New()
	u := usecase.New(r)
	h := httpserver.NewHandler(u)
	router := httpserver.NewRouter(h)

	return &http.Server{
		Addr:    addr,
		Handler: router,
	}
}
