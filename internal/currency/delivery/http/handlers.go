package http

import (
	"currency/internal/currency"
	"currency/pkg/logger"
	"currency/pkg/parser"
	"encoding/json"
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

// GetDates godoc
// @Summary Получение уникальных дат с временем
// @Description Возвращает список уникальных дат с временем, когда были сохранены курсы валют
// @Tags Currency
// @Accept json
// @Produce json
// @Success 200 {array} time.Time
// @Failure 500 {string} string "Internal server error"
// @Router /currency/dates [get]
func (h *currencyHandler) GetDates(w http.ResponseWriter, r *http.Request) {
	dates, err := h.repo.GetDates()
	if err != nil {
		h.logger.Error("GetList error: ", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(dates)
}

// GetCurrency godoc
// @Summary Получение курсов валют по дате и времени
// @Description Возвращает курсы валют по указанной дате и времени
// @Tags Currency
// @Accept json
// @Produce json
// @Param date query string true "Дата и время в формате RFC3339Nano (e.g., 2006-01-02T15:04:05.999999Z)"
// @Success 200 {array} models.Currency
// @Failure 400 {string} string "Invalid date format"
// @Failure 500 {string} string "Internal server error"
// @Router /currency/currency [get]
func (h *currencyHandler) GetCurrency(w http.ResponseWriter, r *http.Request) {
	dateStr := r.URL.Query().Get("date")
	date, err := time.Parse(time.RFC3339Nano, dateStr)
	if err != nil {
		http.Error(w, "Invalid date format. Use RFC3339Nano (e.g., 2006-01-02T15:04:05.999999Z)", http.StatusBadRequest)
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

// GetActual godoc
// @Summary Получение актуальных курсов валют
// @Description Парсит актуальные курсы валют и сохраняет их в базу данных
// @Tags Currency
// @Accept json
// @Produce json
// @Success 200 {array} models.Currency
// @Failure 500 {string} string "Failed to get actual rates"
// @Failure 500 {string} string "Failed to save rates"
// @Router /currency/actual [get]
func (h *currencyHandler) GetActual(w http.ResponseWriter, r *http.Request) {
	currencies, err := parser.ParseCurrencies()
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
