package currency

import (
	"net/http"
)

// Song HTTP Handlers interface
type Handlers interface {
	GetList(w http.ResponseWriter, r *http.Request)
	GetByTime(w http.ResponseWriter, r *http.Request)
	GetActual(w http.ResponseWriter, r *http.Request)
}
