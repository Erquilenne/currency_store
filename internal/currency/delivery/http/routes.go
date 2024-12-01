package http

import (
	"currency/internal/currency"

	"github.com/gorilla/mux"
)

// Map currency routes
func MapCurrencyRoutes(newsGroup *mux.Router, h currency.Handlers) {
	newsGroup.HandleFunc("/dates", h.GetDates).Methods("GET")
	newsGroup.HandleFunc("/currency", h.GetCurrency).Methods("GET")
	newsGroup.HandleFunc("/actual", h.GetActual).Methods("GET")
}
