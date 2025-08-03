package price

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type Service struct{}

func NewService() *Service {
	return &Service{}
}

func (s *Service) GetBitcoinPrice() (string, error) {
	resp, err := http.Get("https://api.coingecko.com/api/v3/simple/price?ids=bitcoin&vs_currencies=usd")
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var data map[string]map[string]float64
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return "", err
	}

	price := data["bitcoin"]["usd"]
	return fmt.Sprintf("%.2f", price), nil
}

func (s *Service) GetEthereumPrice() (string, error) {
	resp, err := http.Get("https://api.coingecko.com/api/v3/simple/price?ids=ethereum&vs_currencies=usd")
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var data map[string]map[string]float64
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return "", err
	}

	price := data["ethereum"]["usd"]
	return fmt.Sprintf("%.2f", price), nil
}
