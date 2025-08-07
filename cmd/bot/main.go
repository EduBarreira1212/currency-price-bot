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

				var lines []string
				for _, c := range telegram.Coins {
					price, err := priceService.GetPrice(c.ID, currency)
					if err != nil {
						lines = append(lines, fmt.Sprintf("%s %s (%s): error: %v", c.Emoji, c.Label, cu, err))
						continue
					}
					lines = append(lines, fmt.Sprintf("%s %s (%s): $%s", c.Emoji, c.Label, cu, price))
				}

				text := strings.Join(lines, "\n")
				bot.SendMessage(tgbotapi.NewMessage(chatID, text))
				bot.UpdateLastSent(chatID)
			}
		}
	}()

	select {}
}
