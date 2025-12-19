package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

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

// ==================================== POST ====================================
func (h *PricesHandler) handlePost(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	file, _, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "FormFile error: "+err.Error(), http.StatusBadRequest)
		return
	}
	defer file.Close()

	records, err := csvzip.ReadCSVFromMultipart(file)
	if err != nil {
		http.Error(w, "CSV read error: "+err.Error(), http.StatusBadRequest)
		return
	}

	if len(records) <= 1 {
		http.Error(w, "CSV contains no data rows", http.StatusBadRequest)
		return
	}

	stats, err := h.db.InsertPrices(r.Context(), records)
	if err != nil {
		http.Error(w, "DB insert error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]int{
		"total_items":      stats.TotalItems,
		"total_categories": stats.TotalCategories,
		"total_price":      stats.TotalPrice,
	})
}

// ==================================== GET ====================================
func (h *PricesHandler) handleGet(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/zip")
	w.Header().Set("Content-Disposition", `attachment; filename="data.zip"`)

	prices, err := h.db.GetPrices(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	rows := make([][]string, 0, len(prices))
	for _, p := range prices {
		rows = append(rows, []string{
			strconv.Itoa(p.ID),
			p.Name,
			p.Category,
			strconv.FormatFloat(p.Price, 'f', -1, 64),
			p.CreateDate.Format("2006-01-02"),
		})
	}

	err = csvzip.WriteCSVToZip(
		w,
		"data.csv",
		[]string{"id", "name", "category", "price", "create_date"},
		rows,
	)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
