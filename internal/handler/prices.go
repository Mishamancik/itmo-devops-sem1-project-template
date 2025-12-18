package handler

import (
	"encoding/json"
	"net/http"
	"archive/zip"
	// "strconv"

	"project_sem/internal/csvzip"
	"project_sem/internal/db"
)

type PricesHandler struct {
	db *db.DB
}

func NewPricesHandler(db *db.DB) *PricesHandler {
	return &PricesHandler{db: db}
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

	file, _, err := r.FormFile("file")
	if err != nil {
		json.NewEncoder(w).Encode(zeroResponse())
		return
	}
	defer file.Close()

	records, err := csvzip.ReadCSVFromMultipart(file)
	if err != nil || len(records) <= 1 {
		json.NewEncoder(w).Encode(zeroResponse())
		return
	}

	stats, err := h.db.InsertPrices(r.Context(), records)
	if err != nil {
		json.NewEncoder(w).Encode(zeroResponse())
		return
	}

	json.NewEncoder(w).Encode(map[string]int{
		"total_items":      stats.TotalItems,
		"total_categories": stats.TotalCategories,
		"total_price":      stats.TotalPrice,
	})
}

func zeroResponse() map[string]int {
	return map[string]int{
		"total_items":      0,
		"total_categories": 0,
		"total_price":      0,
	}
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
