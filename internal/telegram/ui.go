package telegram

import tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

func buildMainKeyboard(subscribed bool) tgbotapi.InlineKeyboardMarkup {
	status := tgbotapi.NewInlineKeyboardButtonData("🔔 Start Updates", "start_updates")
	if subscribed {
		status = tgbotapi.NewInlineKeyboardButtonData("🔕 Stop Updates", "stop_updates")
	}

	btnBTC := tgbotapi.NewInlineKeyboardButtonData("💰 Get BTC", "get_btc")
	btnETH := tgbotapi.NewInlineKeyboardButtonData("🔥 Get ETH", "get_eth")
	btnSOL := tgbotapi.NewInlineKeyboardButtonData("⚡ Get SOL", "get_sol")

	cUSD := tgbotapi.NewInlineKeyboardButtonData("💵 USD", "currency_usd")
	cEUR := tgbotapi.NewInlineKeyboardButtonData("💶 EUR", "currency_eur")
	cBRL := tgbotapi.NewInlineKeyboardButtonData("🇧🇷 BRL", "currency_brl")

	return tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(btnBTC, btnETH, btnSOL),
		tgbotapi.NewInlineKeyboardRow(status),
		tgbotapi.NewInlineKeyboardRow(cUSD, cEUR, cBRL),
	)
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
