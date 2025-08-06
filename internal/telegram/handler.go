package telegram

import (
	"currency-price-bot/internal/price"
	"fmt"
	"strings"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Bot struct {
	api               *tgbotapi.BotAPI
	price             *price.Service
	subscribers       map[int64]bool
	intervals         map[int64]time.Duration
	lastSent          map[int64]time.Time
	preferredCurrency map[int64]string
}

func NewBot(token string, priceService *price.Service) *Bot {
	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		panic(err)
	}

	return &Bot{
		api:               bot,
		price:             priceService,
		subscribers:       make(map[int64]bool),
		intervals:         make(map[int64]time.Duration),
		lastSent:          make(map[int64]time.Time),
		preferredCurrency: make(map[int64]string),
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

	currencyButtons := tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("💵 USD", "currency_usd"),
		tgbotapi.NewInlineKeyboardButtonData("💶 EUR", "currency_eur"),
		tgbotapi.NewInlineKeyboardButtonData("🇧🇷 BRL", "currency_brl"),
	)

	keyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(buttonBTC, buttonETH, buttonSOL),
		tgbotapi.NewInlineKeyboardRow(statusButton),
		currencyButtons,
	)

	text := "Welcome! Choose a coin or manage auto updates:"
	message := tgbotapi.NewMessage(msg.Chat.ID, text)
	message.ReplyMarkup = keyboard

	b.api.Send(message)
}

func (b *Bot) handleStartUpdatesCommand(msg *tgbotapi.Message) {
	var statusButton tgbotapi.InlineKeyboardButton

	if b.subscribers[msg.Chat.ID] {
		statusButton = tgbotapi.NewInlineKeyboardButtonData("🔕 Stop Updates", "stop_updates")
	} else {
		statusButton = tgbotapi.NewInlineKeyboardButtonData("🔔 Start Updates", "start_updates")
	}

	intervalButtons := tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("1️⃣ 1 min", "interval_1"),
		tgbotapi.NewInlineKeyboardButtonData("5️⃣ 5 min", "interval_5"),
		tgbotapi.NewInlineKeyboardButtonData("🔟 10 min", "interval_10"),
	)

	keyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(statusButton),
		intervalButtons,
	)

	text := "Choose an interval:"
	message := tgbotapi.NewMessage(msg.Chat.ID, text)
	message.ReplyMarkup = keyboard

	b.api.Send(message)
}

func (b *Bot) sendPrice(chatID int64, callbackID, coin string) {
	b.api.Request(tgbotapi.NewCallback(callbackID, "Fetching "+coin+" price..."))

	currency := b.getCurrency(chatID)

	price, err := b.price.GetPrice(coin, currency)
	text := fmt.Sprintf("💵 %s (%s): %s", coin, strings.ToUpper(currency), price)

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
	case "currency_usd":
		b.setCurrency(cb.Message.Chat.ID, "usd")
		b.api.Send(tgbotapi.NewMessage(cb.Message.Chat.ID, "💵 Currency set to USD."))
	case "currency_eur":
		b.setCurrency(cb.Message.Chat.ID, "eur")
		b.api.Send(tgbotapi.NewMessage(cb.Message.Chat.ID, "💶 Currency set to EUR."))
	case "currency_brl":
		b.setCurrency(cb.Message.Chat.ID, "brl")
		b.api.Send(tgbotapi.NewMessage(cb.Message.Chat.ID, "🇧🇷 Currency set to BRL."))
	case "start_updates":
		b.subscribers[cb.Message.Chat.ID] = true
		b.api.Send(tgbotapi.NewMessage(cb.Message.Chat.ID, "🔔 Auto-updates enabled!"))
		b.handleStartUpdatesCommand(cb.Message)
	case "stop_updates":
		delete(b.subscribers, cb.Message.Chat.ID)
		b.api.Send(tgbotapi.NewMessage(cb.Message.Chat.ID, "🔕 Auto-updates disabled."))
		b.handleStartCommand(cb.Message)
	case "interval_1":
		b.setInterval(cb.Message.Chat.ID, 1*time.Minute)
		b.api.Send(tgbotapi.NewMessage(cb.Message.Chat.ID, "🕐 Interval set to 1 minute."))
	case "interval_5":
		b.setInterval(cb.Message.Chat.ID, 5*time.Minute)
		b.api.Send(tgbotapi.NewMessage(cb.Message.Chat.ID, "🕔 Interval set to 5 minutes."))
	case "interval_10":
		b.setInterval(cb.Message.Chat.ID, 10*time.Minute)
		b.api.Send(tgbotapi.NewMessage(cb.Message.Chat.ID, "🔟 Interval set to 10 minutes."))
	}
}

func (b *Bot) Subscribers() map[int64]bool {
	return b.subscribers
}

func (b *Bot) SendMessage(msg tgbotapi.MessageConfig) {
	b.api.Send(msg)
}

func (b *Bot) setInterval(chatID int64, interval time.Duration) {
	b.intervals[chatID] = interval
	if _, ok := b.lastSent[chatID]; !ok {
		b.lastSent[chatID] = time.Time{}
	}
}

func (b *Bot) GetInterval(chatID int64) time.Duration {
	if interval, ok := b.intervals[chatID]; ok {
		return interval
	}
	return 5 * time.Minute
}

func (b *Bot) GetLastSent(chatID int64) time.Time {
	if last, ok := b.lastSent[chatID]; ok {
		return last
	}
	return time.Time{}
}

func (b *Bot) UpdateLastSent(chatID int64) {
	b.lastSent[chatID] = time.Now()
}

func (b *Bot) setCurrency(chatID int64, currency string) {
	b.preferredCurrency[chatID] = currency
}

func (b *Bot) getCurrency(chatID int64) string {
	if c, ok := b.preferredCurrency[chatID]; ok {
		return c
	}
	return "usd"
}

func (b *Bot) GetCurrency(chatID int64) string {
	return b.getCurrency(chatID)
}
