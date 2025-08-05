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

func (s *Service) GetPrice(coin string) (string, error) {
	url := fmt.Sprintf("https://api.coingecko.com/api/v3/simple/price?ids=%s&vs_currencies=usd", coin)
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var data map[string]map[string]float64
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return "", err
	}

	price := data[coin]["usd"]
	return fmt.Sprintf("%.2f", price), nil
}
