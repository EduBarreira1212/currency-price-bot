package telegram

import (
	"currency-price-bot/internal/price"
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Bot struct {
	api   *tgbotapi.BotAPI
	price *price.Service

	state botState
}

func NewBot(token string, priceService *price.Service) *Bot {
	api, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		panic(err)
	}
	return &Bot{
		api:   api,
		price: priceService,
		state: newBotState(),
	}
}

func (b *Bot) Start() {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	updates := b.api.GetUpdatesChan(u)

	for update := range updates {
		b.handleUpdate(update)
	}
}

func (b *Bot) SendMessage(msg tgbotapi.MessageConfig) {
	if _, err := b.api.Send(msg); err != nil {
		log.Printf("send error: %v", err)
	}
}

func (b *Bot) SnapshotSubscribers() []int64              { return b.state.snapshotSubscribers() }
func (b *Bot) GetInterval(chatID int64) (d DurationLike) { return b.state.getInterval(chatID) }
func (b *Bot) GetLastSent(chatID int64) (t TimeLike)     { return b.state.getLastSent(chatID) }
func (b *Bot) UpdateLastSent(chatID int64)               { b.state.updateLastSent(chatID) }
func (b *Bot) GetCurrency(chatID int64) string           { return b.state.getCurrency(chatID) }
