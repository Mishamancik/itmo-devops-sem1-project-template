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

	// Health endpoint (без проверки БД)
	r.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("ok"))
	}).Methods(http.MethodGet)

	// Подключение к БД
	database, err := db.New()
	if err != nil {
		log.Fatal(err)
	}

	// Хендлеры API
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
