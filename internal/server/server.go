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

	// Здесь создаём подключение к БД
	database, err := db.New()
	if err != nil {
		log.Fatal(err)
	}

	// Здесь передаём БД в handler
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
