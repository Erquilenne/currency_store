package server

import (
	currencyHttp "currency/internal/currency/delivery/http"
	"currency/internal/currency/repository"

	"github.com/gorilla/mux"
)

// MapHandlers Map Server Handlers
func (s *Server) MapHandlers(router *mux.Router) error {
	currencyRepo := repository.NewCurrencyRepository(s.db)

	currencyHandlers := currencyHttp.NewCurrencyHandler(currencyRepo, s.logger)

	apiRouter := router.PathPrefix("/api/v1").Subrouter()

	currencyGroup := apiRouter.PathPrefix("/currency").Subrouter()
	currencyHttp.MapCurrencyRoutes(currencyGroup, currencyHandlers)

	return nil
}
