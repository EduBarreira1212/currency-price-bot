package main

import (
	"currency-price-bot/internal/price"
	"currency-price-bot/internal/telegram"
	"fmt"
	"log"
	"os"
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
			time.Sleep(30 * time.Second)

			for chatID := range bot.Subscribers() {
				btcPrice, _ := priceService.GetPrice("bitcoin")
				ethPrice, _ := priceService.GetPrice("ethereum")
				solPrice, _ := priceService.GetPrice("solana")

				text := fmt.Sprintf("📈 BTC: $%s\n📉 ETH: $%s\n📈 SOL: $%s", btcPrice, ethPrice, solPrice)

				msg := tgbotapi.NewMessage(chatID, text)
				bot.SendMessage(msg)
			}
		}
	}()

	select {}
}
