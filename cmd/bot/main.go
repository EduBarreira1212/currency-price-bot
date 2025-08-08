package main

import (
	"currency-price-bot/internal/price"
	"currency-price-bot/internal/telegram"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	token := os.Getenv("TELEGRAM_BOT_TOKEN")
	if token == "" {
		log.Fatal("TELEGRAM_BOT_TOKEN is not set")
	}

	priceService := price.NewService()
	bot := telegram.NewBot(token, priceService)
	log.Println("Bot started successfully")

	go bot.Start()

	coinIDs := make([]string, 0, len(telegram.Coins))
	for _, c := range telegram.Coins {
		coinIDs = append(coinIDs, c.ID)
	}

	ticker := time.NewTicker(10 * time.Second)
	go func() {
		for range ticker.C {
			for _, chatID := range bot.SnapshotSubscribers() {
				interval := bot.GetInterval(chatID)
				last := bot.GetLastSent(chatID)
				if time.Since(last) < interval {
					continue
				}

				currency := bot.GetCurrency(chatID)
				cu := strings.ToUpper(currency)

				prices, err := priceService.GetPrices(coinIDs, currency)
				if err != nil {
					log.Printf("batch GetPrices error: %v (falling back to single requests)", err)
					prices = make(map[string]string, len(coinIDs))
					for _, id := range coinIDs {
						p, e := priceService.GetPrice(id, currency)
						if e != nil {
							prices[id] = "__ERR__:" + e.Error()
						} else {
							prices[id] = p
						}
					}
				}

				var lines []string
				for _, c := range telegram.Coins {
					val, ok := prices[c.ID]
					if !ok {
						lines = append(lines, fmt.Sprintf("%s %s (%s): unavailable", c.Emoji, c.Label, cu))
						continue
					}
					if strings.HasPrefix(val, "__ERR__:") {
						lines = append(lines, fmt.Sprintf("%s %s (%s): error: %s", c.Emoji, c.Label, cu, strings.TrimPrefix(val, "__ERR__:")))
						continue
					}
					lines = append(lines, fmt.Sprintf("%s %s (%s): $%s", c.Emoji, c.Label, cu, val))
				}

				text := strings.Join(lines, "\n")
				bot.SendMessage(tgbotapi.NewMessage(chatID, text))
				bot.UpdateLastSent(chatID)
			}
		}
	}()

	select {}
}
