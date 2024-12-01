package http

import (
	"currency/internal/currency"

	"github.com/gorilla/mux"
)

// Map currency routes
func MapCurrencyRoutes(newsGroup *mux.Router, h currency.Handlers) {
	newsGroup.HandleFunc("/list", h.GetList).Methods("GET")
	newsGroup.HandleFunc("/", h.GetByTime).Methods("GET")
	newsGroup.HandleFunc("/actual", h.GetActual).Methods("GET")
}
