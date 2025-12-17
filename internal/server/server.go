package server

import (
	"net/http"

	"github.com/gorilla/mux"
	"project_sem/internal/handler"
)

type Server struct {
	router *mux.Router
}

func New() *Server {
	r := mux.NewRouter()

	h := handler.NewPricesHandler()

	r.HandleFunc("/api/v0/prices", h.HandlePrices).
		Methods(http.MethodPost, http.MethodGet)

	return &Server{
		router: r,
	}
}

func (s *Server) Start(addr string) error {
	return http.ListenAndServe(addr, s.router)
}
