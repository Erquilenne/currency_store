package repository

import (
	"currency/internal/models"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
)

type currencyRepo struct {
	db *sqlx.DB
}

func NewCurrencyRepository(db *sqlx.DB) *currencyRepo {
	return &currencyRepo{db: db}
}

func (r *currencyRepo) GetDates() ([]time.Time, error) {
	var dates []time.Time
	err := r.db.Select(&dates, getDates)
	fmt.Println(dates)
	return dates, err
}

func (r *currencyRepo) GetByDate(date time.Time) ([]models.Currency, error) {
	var currencies []models.Currency
	err := r.db.Select(&currencies, getByDate, date)
	return currencies, err
}

func (r *currencyRepo) SaveCurrencies(currencies []models.Currency) error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	for _, curr := range currencies {
		_, err = tx.Exec(save, curr.Currency, curr.Type, curr.Value)
		if err != nil {
			return err
		}
	}

	return tx.Commit()
}
