package currency

import (
	"currency/internal/models"
	"time"
)

// Repository interface
type Repository interface {
	GetDates() ([]time.Time, error)
	GetByDate(date time.Time) ([]models.Currency, error)
	SaveCurrencies(currencies []models.Currency) error
}
