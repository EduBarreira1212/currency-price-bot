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

				btcPrice, err1 := priceService.GetPrice("bitcoin", currency)
				ethPrice, err2 := priceService.GetPrice("ethereum", currency)
				solPrice, err3 := priceService.GetPrice("solana", currency)

				var text string
				if err1 != nil || err2 != nil || err3 != nil {
					text = "⚠️ Failed fetching some prices:\n"
					if err1 != nil {
						text += "BTC error: " + err1.Error() + "\n"
					} else {
						text += fmt.Sprintf("📈 BTC (%s): $%s\n", cu, btcPrice)
					}
					if err2 != nil {
						text += "ETH error: " + err2.Error() + "\n"
					} else {
						text += fmt.Sprintf("📉 ETH (%s): $%s\n", cu, ethPrice)
					}
					if err3 != nil {
						text += "SOL error: " + err3.Error()
					} else {
						text += fmt.Sprintf("⚡ SOL (%s): $%s", cu, solPrice)
					}
				} else {
					text = fmt.Sprintf(
						"📈 BTC (%s): $%s\n📉 ETH (%s): $%s\n⚡ SOL (%s): $%s",
						cu, btcPrice, cu, ethPrice, cu, solPrice,
					)
				}

				bot.SendMessage(tgbotapi.NewMessage(chatID, text))
				bot.UpdateLastSent(chatID)
			}
		}
	}()

	select {}
}
