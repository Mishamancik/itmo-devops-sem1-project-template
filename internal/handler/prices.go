package handler

import (
	"encoding/json"
	"net/http"
	"archive/zip"
)

type PricesHandler struct{}

func NewPricesHandler() *PricesHandler {
	return &PricesHandler{}
}

func (h *PricesHandler) HandlePrices(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		h.handlePost(w, r)
	case http.MethodGet:
		h.handleGet(w, r)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func (h *PricesHandler) handlePost(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Пока заглушка
	resp := map[string]int{
		"total_items":      0,
		"total_categories": 0,
		"total_price":      0,
	}

	_ = json.NewEncoder(w).Encode(resp)
}

func (h *PricesHandler) handleGet(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/zip")
	w.Header().Set("Content-Disposition", `attachment; filename="data.zip"`)

	zipWriter := zip.NewWriter(w)

	fileWriter, err := zipWriter.Create("data.csv")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Пока только заголовки
	_, err = fileWriter.Write([]byte("id,name,category,price,create_date\n"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := zipWriter.Close(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
