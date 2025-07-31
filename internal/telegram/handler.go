package telegram

import (
	"currency-price-bot/internal/price"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Bot struct {
	api   *tgbotapi.BotAPI
	price *price.Service
}

func NewBot(token string, priceService *price.Service) *Bot {
	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		panic(err)
	}

	return &Bot{api: bot, price: priceService}
}

func (b *Bot) Start() {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	updates := b.api.GetUpdatesChan(u)

	for update := range updates {
		if update.Message == nil {
			continue
		}

		if update.Message.Text == "/btc" {
			btcPrice, err := b.price.GetBitcoinPrice()
			msgText := "Bitcoin price: $" + btcPrice
			if err != nil {
				msgText = "Failed to fetch price: " + err.Error()
			}

			msg := tgbotapi.NewMessage(update.Message.Chat.ID, msgText)
			b.api.Send(msg)
		}
	}
}
