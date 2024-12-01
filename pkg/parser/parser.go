package parser

import (
	"context"
	"currency/internal/models"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

const (
	alfabankAPIBaseURL = "https://alfabank.ru/api/v1/scrooge/currencies/alfa-rates"
	timeout            = 30 * time.Second // увеличиваем таймаут
)

type AlfaBankResponse struct {
	Data []struct {
		CurrencyCode     string `json:"currencyCode"`
		RateByClientType []struct {
			ClientType  string `json:"clientType"`
			RatesByType []struct {
				RateType       string `json:"rateType"`
				LastActualRate struct {
					Buy struct {
						OriginalValue float64 `json:"originalValue"`
					} `json:"buy"`
					Sell struct {
						OriginalValue float64 `json:"originalValue"`
					} `json:"sell"`
				} `json:"lastActualRate"`
			} `json:"ratesByType"`
		} `json:"rateByClientType"`
	} `json:"data"`
}

func ParseCurrencies() ([]models.Currency, error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	// Получаем текущую дату и время в нужном формате
	currentTime := time.Now().Format("2006-01-02T15:04:05+03:00")

	// Конструируем URL запроса с текущей датой и отфильтрованными валютами
	requestURL := fmt.Sprintf("%s?clientType.eq=standardCC&currencyCode.in=USD,EUR&date.lte=%s&lastActualForDate.eq=true&rateType.in=rateCass", alfabankAPIBaseURL, currentTime)

	req, err := http.NewRequestWithContext(ctx, "GET", requestURL, nil)
	if err != nil {
		return nil, fmt.Errorf("не удалось создать запрос: %w", err)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("не удалось выполнить запрос: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("получен некорректный статус код: %d", resp.StatusCode)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("не удалось прочитать тело ответа: %w", err)
	}

	var response AlfaBankResponse
	err = json.Unmarshal(body, &response)
	if err != nil {
		return nil, fmt.Errorf("не удалось распарсить JSON: %w", err)
	}

	var currencies []models.Currency
	for _, item := range response.Data {
		currencyCode := strings.ToLower(item.CurrencyCode)
		for _, rateByClientType := range item.RateByClientType {
			for _, ratesByType := range rateByClientType.RatesByType {
				currencies = append(currencies, models.Currency{
					Currency: currencyCode,
					Type:     "buy",
					Value:    ratesByType.LastActualRate.Buy.OriginalValue,
				})
				currencies = append(currencies, models.Currency{
					Currency: currencyCode,
					Type:     "sell",
					Value:    ratesByType.LastActualRate.Sell.OriginalValue,
				})
			}
		}
	}

	return currencies, nil
}
