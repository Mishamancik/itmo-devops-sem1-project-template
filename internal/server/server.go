package server

import (
	"net/http"

	"github.com/gorilla/mux"

	"project_sem/internal/db"
	"project_sem/internal/handler"
)

type Server struct {
	router *mux.Router
}

func New(database *db.DB) (*Server, error) {
	r := mux.NewRouter()

	r.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ok"))
	}).Methods(http.MethodGet)

	h := handler.NewPricesHandler(database)

	r.HandleFunc("/api/v0/prices", h.HandlePrices).
		Methods(http.MethodPost, http.MethodGet)

	return &Server{
		router: r,
	}, nil
}

func (s *Server) Start(addr string) error {
	return http.ListenAndServe(addr, s.router)
}
