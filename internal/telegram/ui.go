package telegram

import tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

func buildMainKeyboard(subscribed bool) tgbotapi.InlineKeyboardMarkup {
	status := tgbotapi.NewInlineKeyboardButtonData("🔔 Start Updates", "start_updates")
	if subscribed {
		status = tgbotapi.NewInlineKeyboardButtonData("🔕 Stop Updates", "stop_updates")
	}

	rows := [][]tgbotapi.InlineKeyboardButton{}
	row := []tgbotapi.InlineKeyboardButton{}
	for i, c := range Coins {
		btn := tgbotapi.NewInlineKeyboardButtonData(c.Emoji+" Get "+c.Label, "get_"+c.ID)
		row = append(row, btn)
		if (i+1)%3 == 0 {
			rows = append(rows, row)
			row = []tgbotapi.InlineKeyboardButton{}
		}
	}
	if len(row) > 0 {
		rows = append(rows, row)
	}

	cUSD := tgbotapi.NewInlineKeyboardButtonData("💵 USD", "currency_usd")
	cEUR := tgbotapi.NewInlineKeyboardButtonData("💶 EUR", "currency_eur")
	cBRL := tgbotapi.NewInlineKeyboardButtonData("🇧🇷 BRL", "currency_brl")

	rows = append(rows,
		tgbotapi.NewInlineKeyboardRow(status),
		tgbotapi.NewInlineKeyboardRow(cUSD, cEUR, cBRL),
	)

	return tgbotapi.NewInlineKeyboardMarkup(rows...)
}

func buildIntervalKeyboard(subscribed bool) tgbotapi.InlineKeyboardMarkup {
	status := tgbotapi.NewInlineKeyboardButtonData("🔔 Start Updates", "start_updates")
	if subscribed {
		status = tgbotapi.NewInlineKeyboardButtonData("🔕 Stop Updates", "stop_updates")
	}
	i1 := tgbotapi.NewInlineKeyboardButtonData("1️⃣ 1 min", "interval_1")
	i5 := tgbotapi.NewInlineKeyboardButtonData("5️⃣ 5 min", "interval_5")
	i10 := tgbotapi.NewInlineKeyboardButtonData("🔟 10 min", "interval_10")

	return tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(status),
		tgbotapi.NewInlineKeyboardRow(i1, i5, i10),
	)
}
