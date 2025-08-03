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
		if update.Message != nil && update.Message.Text == "/start" {
			b.handleStartCommand(update.Message)
			continue
		}

		if update.CallbackQuery != nil {
			b.handleCallback(update.CallbackQuery)
			continue
		}
	}
}

func (b *Bot) handleStartCommand(msg *tgbotapi.Message) {
	btcButton := tgbotapi.NewInlineKeyboardButtonData("💰 Get Bitcoin Price", "get_btc")
	ethButton := tgbotapi.NewInlineKeyboardButtonData("💰 Get Ethereum Price", "get_eth")
	keyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(btcButton, ethButton),
	)

	text := "Welcome! Click the button below to get the price of the choosen coin:"
	message := tgbotapi.NewMessage(msg.Chat.ID, text)
	message.ReplyMarkup = keyboard

	b.api.Send(message)
}

func (b *Bot) handleCallback(cb *tgbotapi.CallbackQuery) {
	switch cb.Data {
	case "get_btc":
		b.api.Request(tgbotapi.NewCallback(cb.ID, "Fetching BTC price..."))

		price, err := b.price.GetBitcoinPrice()
		msgText := "💵 BTC: $" + price
		if err != nil {
			msgText = "❌ Error fetching BTC price: " + err.Error()
		}

		msg := tgbotapi.NewMessage(cb.Message.Chat.ID, msgText)
		b.api.Send(msg)
	case "get_eth":
		b.api.Request(tgbotapi.NewCallback(cb.ID, "Fetching ETH price..."))

		price, err := b.price.GetEthereumPrice()
		msgText := "💵 ETH: $" + price
		if err != nil {
			msgText = "❌ Error fetching ETH price: " + err.Error()
		}

		msg := tgbotapi.NewMessage(cb.Message.Chat.ID, msgText)
		b.api.Send(msg)
	}
}
