package main

import (
	"currency-price-bot/internal/price"
	"currency-price-bot/internal/telegram"
	"log"
	"os"

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
	bot.Start()
}
