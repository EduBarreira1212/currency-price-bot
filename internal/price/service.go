package price

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"
)

type Service struct {
	http *http.Client
}

func NewService() *Service {
	return &Service{
		http: &http.Client{Timeout: 10 * time.Second},
	}
}

func (s *Service) GetPrice(coinID, currency string) (string, error) {
	prices, err := s.GetPrices([]string{coinID}, currency)
	if err != nil {
		return "", err
	}
	p, ok := prices[coinID]
	if !ok {
		return "", fmt.Errorf("no price for %s", coinID)
	}
	return p, nil
}

func (s *Service) GetPrices(coinIDs []string, currency string) (map[string]string, error) {
	if len(coinIDs) == 0 {
		return nil, errors.New("no coin ids")
	}

	ids := ""
	for i, id := range coinIDs {
		if i > 0 {
			ids += ","
		}
		ids += id
	}

	url := fmt.Sprintf(
		"https://api.coingecko.com/api/v3/simple/price?ids=%s&vs_currencies=%s",
		ids, currency,
	)

	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("User-Agent", "currency-price-bot/1.0 (+https://github.com/)")

	resp, err := s.http.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("coingecko status %d", resp.StatusCode)
	}

	dec := json.NewDecoder(resp.Body)
	dec.UseNumber()

	var raw map[string]map[string]interface{}
	if err := dec.Decode(&raw); err != nil {
		return nil, err
	}

	out := make(map[string]string, len(coinIDs))
	for _, id := range coinIDs {
		curMap, ok := raw[id]
		if !ok {
			continue
		}
		val, ok := curMap[currency]
		if !ok {
			continue
		}

		f, err := toFloat(val)
		if err != nil {
			return nil, fmt.Errorf("%s (%s): %w", id, currency, err)
		}
		out[id] = fmt.Sprintf("%.2f", f)
	}

	return out, nil
}

func toFloat(v interface{}) (float64, error) {
	switch t := v.(type) {
	case json.Number:
		return t.Float64()
	case float64:
		return t, nil
	case string:
		if t == "NaN" || t == "" {
			return 0, fmt.Errorf("non-numeric string %q", t)
		}
		f, err := strconv.ParseFloat(t, 64)
		if err != nil {
			return 0, fmt.Errorf("cannot parse %q: %w", t, err)
		}
		return f, nil
	default:
		return 0, fmt.Errorf("unsupported JSON type %T", v)
	}
}
