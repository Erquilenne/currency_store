package models

type Currency struct {
	Currency string  `json:"currency"` // usd или eur
	Type     string  `json:"type"`     // buy или sell
	Value    float64 `json:"value"`
}
