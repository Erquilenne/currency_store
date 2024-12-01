package models

type Currency struct {
	Currency string  // USD, EUR и т.д.
	Type     string  // buy или sell
	Value    float64 // значение курса
}
