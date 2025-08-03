package telegram

import (
	"currency-price-bot/internal/price"
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Bot struct {
	api         *tgbotapi.BotAPI
	price       *price.Service
	subscribers map[int64]bool
}

func NewBot(token string, priceService *price.Service) *Bot {
	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		panic(err)
	}

	return &Bot{
		api:         bot,
		price:       priceService,
		subscribers: make(map[int64]bool),
	}
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
	var statusButton tgbotapi.InlineKeyboardButton

	if b.subscribers[msg.Chat.ID] {
		statusButton = tgbotapi.NewInlineKeyboardButtonData("🔕 Stop Updates", "stop_updates")
	} else {
		statusButton = tgbotapi.NewInlineKeyboardButtonData("🔔 Start Updates", "start_updates")
	}

	buttonBTC := tgbotapi.NewInlineKeyboardButtonData("💰 Get BTC", "get_btc")
	buttonETH := tgbotapi.NewInlineKeyboardButtonData("🔥 Get ETH", "get_eth")
	buttonSOL := tgbotapi.NewInlineKeyboardButtonData("🔥 Get SOL", "get_sol")

	keyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(buttonBTC, buttonETH, buttonSOL),
		tgbotapi.NewInlineKeyboardRow(statusButton),
	)

	text := "Welcome! Choose a coin or manage auto updates:"
	message := tgbotapi.NewMessage(msg.Chat.ID, text)
	message.ReplyMarkup = keyboard

	b.api.Send(message)
}

func (b *Bot) sendPrice(chatID int64, callbackID, coin string) {
	b.api.Request(tgbotapi.NewCallback(callbackID, "Fetching "+coin+" price..."))

	price, err := b.price.GetPrice(coin)
	text := fmt.Sprintf("💵 %s: $%s", coin, price)
	if err != nil {
		text = "❌ Error fetching price: " + err.Error()
	}

	msg := tgbotapi.NewMessage(chatID, text)
	b.api.Send(msg)
}

func (b *Bot) handleCallback(cb *tgbotapi.CallbackQuery) {
	switch cb.Data {
	case "get_btc":
		b.sendPrice(cb.Message.Chat.ID, cb.ID, "bitcoin")
	case "get_eth":
		b.sendPrice(cb.Message.Chat.ID, cb.ID, "ethereum")
	case "get_sol":
		b.sendPrice(cb.Message.Chat.ID, cb.ID, "solana")
	case "start_updates":
		b.subscribers[cb.Message.Chat.ID] = true
		b.api.Send(tgbotapi.NewMessage(cb.Message.Chat.ID, "🔔 Auto-updates enabled!"))
		b.handleStartCommand(cb.Message)

	case "stop_updates":
		delete(b.subscribers, cb.Message.Chat.ID)
		b.api.Send(tgbotapi.NewMessage(cb.Message.Chat.ID, "🔕 Auto-updates disabled."))
		b.handleStartCommand(cb.Message)

	}
}

func (b *Bot) Subscribers() map[int64]bool {
	return b.subscribers
}

func (b *Bot) SendMessage(msg tgbotapi.MessageConfig) {
	b.api.Send(msg)
}
