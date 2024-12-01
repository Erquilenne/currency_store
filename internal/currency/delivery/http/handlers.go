package http

import (
	"currency/internal/currency"
	"currency/pkg/logger"
	"currency/pkg/parser"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type currencyHandler struct {
	repo   currency.Repository
	logger logger.Logger
}

func NewCurrencyHandler(repo currency.Repository, logger logger.Logger) currency.Handlers {
	return &currencyHandler{
		repo:   repo,
		logger: logger,
	}
}

func (h *currencyHandler) GetList(w http.ResponseWriter, r *http.Request) {
	dates, err := h.repo.GetDates()
	if err != nil {
		h.logger.Error("GetList error: ", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(dates)
}

func (h *currencyHandler) GetByTime(w http.ResponseWriter, r *http.Request) {
	dateStr := r.URL.Query().Get("date")
	date, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		http.Error(w, "Invalid date format. Use YYYY-MM-DD", http.StatusBadRequest)
		return
	}

	currencies, err := h.repo.GetByDate(date)
	if err != nil {
		h.logger.Error("GetByTime error: ", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(currencies)
}

func (h *currencyHandler) GetActual(w http.ResponseWriter, r *http.Request) {
	currencies, err := parser.ParseCurrencies()
	fmt.Println(currencies)
	if err != nil {
		h.logger.Error("Parser error: ", err)
		http.Error(w, "Failed to get actual rates", http.StatusInternalServerError)
		return
	}

	// Сохраняем полученные данные в базу
	if err := h.repo.SaveCurrencies(currencies); err != nil {
		h.logger.Error("Save currencies error: ", err)
		http.Error(w, "Failed to save rates", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(currencies)
}
