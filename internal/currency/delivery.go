package currency

import (
	"net/http"
)

// Song HTTP Handlers interface
type Handlers interface {
	GetDates(w http.ResponseWriter, r *http.Request)
	GetCurrency(w http.ResponseWriter, r *http.Request)
	GetActual(w http.ResponseWriter, r *http.Request)
}
