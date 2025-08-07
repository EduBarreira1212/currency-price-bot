package telegram

import (
	"fmt"
	"strings"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (b *Bot) handleUpdate(u tgbotapi.Update) {
	switch {
	case u.Message != nil:
		b.handleMessage(u.Message)
	case u.CallbackQuery != nil:
		b.handleCallback(u.CallbackQuery)
	}
}

func (b *Bot) handleMessage(msg *tgbotapi.Message) {
	if msg.Text == "/start" {
		b.state.ensureUser(msg.Chat.ID)
		kb := buildMainKeyboard(b.state.isSubscribed(msg.Chat.ID))
		m := tgbotapi.NewMessage(msg.Chat.ID, "Welcome! Choose a coin or manage auto updates:")
		m.ReplyMarkup = kb
		b.api.Send(m)
		return
	}
}

func (b *Bot) handleCallback(cb *tgbotapi.CallbackQuery) {
	data := cb.Data

	if strings.HasPrefix(data, "get_") {
		coinID := strings.TrimPrefix(data, "get_")
		if c, ok := FindCoin(coinID); ok {
			handleGetPrice(c.ID)(b, cb)
			return
		}

		b.api.Request(tgbotapi.NewCallback(cb.ID, "Unknown asset"))
		return
	}

	if h, ok := callbackHandlers[data]; ok {
		h(b, cb)
	}
}

var callbackHandlers = map[string]func(*Bot, *tgbotapi.CallbackQuery){
	"currency_usd": handleSetCurrency("usd"),
	"currency_eur": handleSetCurrency("eur"),
	"currency_brl": handleSetCurrency("brl"),
	"start_updates": func(b *Bot, cb *tgbotapi.CallbackQuery) {
		chatID := cb.Message.Chat.ID
		b.state.setSubscribed(chatID, true)
		b.api.Send(tgbotapi.NewMessage(chatID, "🔔 Auto-updates enabled!"))
		kb := buildIntervalKeyboard(true)
		m := tgbotapi.NewMessage(chatID, "Choose an interval:")
		m.ReplyMarkup = kb
		b.api.Send(m)
	},
	"stop_updates": func(b *Bot, cb *tgbotapi.CallbackQuery) {
		chatID := cb.Message.Chat.ID
		b.state.setSubscribed(chatID, false)
		b.api.Send(tgbotapi.NewMessage(chatID, "🔕 Auto-updates disabled."))
		kb := buildMainKeyboard(false)
		m := tgbotapi.NewMessage(chatID, "Updates stopped. You can enable them again below:")
		m.ReplyMarkup = kb
		b.api.Send(m)
	},
	"interval_1":  handleSetInterval(1 * time.Minute),
	"interval_5":  handleSetInterval(5 * time.Minute),
	"interval_10": handleSetInterval(10 * time.Minute),
}

func handleGetPrice(coinID string) func(*Bot, *tgbotapi.CallbackQuery) {
	return func(b *Bot, cb *tgbotapi.CallbackQuery) {
		chatID := cb.Message.Chat.ID
		b.api.Request(tgbotapi.NewCallback(cb.ID, "Fetching "+coinID+" price..."))
		currency := b.state.getCurrency(chatID)

		p, err := b.price.GetPrice(coinID, currency)
		label := coinID
		if c, ok := FindCoin(coinID); ok {
			label = c.Label
		}
		text := fmt.Sprintf("💵 %s (%s): %s", label, strings.ToUpper(currency), p)
		if err != nil {
			text = "❌ Error fetching price: " + err.Error()
		}
		b.api.Send(tgbotapi.NewMessage(chatID, text))
	}
}

func handleSetCurrency(curr string) func(*Bot, *tgbotapi.CallbackQuery) {
	return func(b *Bot, cb *tgbotapi.CallbackQuery) {
		chatID := cb.Message.Chat.ID
		b.state.setCurrency(chatID, curr)
		label := map[string]string{"usd": "💵 USD", "eur": "💶 EUR", "brl": "🇧🇷 BRL"}[curr]
		if label == "" {
			label = strings.ToUpper(curr)
		}
		b.api.Send(tgbotapi.NewMessage(chatID, fmt.Sprintf("%s selected.", label)))
	}
}

func handleSetInterval(d time.Duration) func(*Bot, *tgbotapi.CallbackQuery) {
	return func(b *Bot, cb *tgbotapi.CallbackQuery) {
		chatID := cb.Message.Chat.ID
		b.state.setInterval(chatID, d)
		b.api.Send(tgbotapi.NewMessage(chatID, fmt.Sprintf("⏱ Interval set to %s.", d)))
	}
}
