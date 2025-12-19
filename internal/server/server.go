package server

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"

	"project_sem/internal/db"
	"project_sem/internal/handler"
)

type Server struct {
	router *mux.Router
}

func New() *Server {
	r := mux.NewRouter()

	r.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("ok"))
	}).Methods(http.MethodGet)

	database, err := db.New()
	if err != nil {
		log.Fatal(err)
	}

	h := handler.NewPricesHandler(database)

	r.HandleFunc("/api/v0/prices", h.HandlePrices).
		Methods(http.MethodPost, http.MethodGet)

	return &Server{
		router: r,
	}
}

func (s *Server) Start(addr string) error {
	return http.ListenAndServe(addr, s.router)
}
