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
	err := godotenv.Load()
	if err != nil {
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
	go func() {
		for {
			time.Sleep(10 * time.Second)

			for chatID := range bot.Subscribers() {
				interval := bot.GetInterval(chatID)
				last := bot.GetLastSent(chatID)

				if time.Since(last) < interval {
					continue
				}

				currency := bot.GetCurrency(chatID)

				btcPrice, _ := priceService.GetPrice("bitcoin", currency)
				ethPrice, _ := priceService.GetPrice("ethereum", currency)
				solPrice, _ := priceService.GetPrice("solana", currency)

				text := fmt.Sprintf("📈 BTC (%s): $%s\n📉 ETH (%s): $%s\n📉 SOL (%s): $%s",
					strings.ToUpper(currency), btcPrice,
					strings.ToUpper(currency), ethPrice,
					strings.ToUpper(currency), solPrice)

				msg := tgbotapi.NewMessage(chatID, text)
				bot.SendMessage(msg)
				bot.UpdateLastSent(chatID)
			}
		}
	}()

	select {}
}
