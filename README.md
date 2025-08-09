# 💹 Crypto Price Telegram Bot

A **Golang Telegram bot** that fetches cryptocurrency prices from [CoinGecko](https://www.coingecko.com/) and sends them directly to you in Telegram.  
Supports multiple coins, custom update intervals, and different fiat currencies.

## ✨ Features

- 📊 **Multiple coins** (BTC, ETH, SOL, ADA, DOGE… easy to add more)
- 🔔 **Start/Stop auto-updates** with one click
- ⏱ **Custom update intervals** (1, 5, 10 minutes)
- 💱 **Preferred currency** (USD, EUR, BRL)
- 🛡 **Safe API calls** with error handling for unexpected JSON
- 🧩 **Modular architecture** for easy maintenance and scaling

---

## 📂 Project Structure

```
.
├── cmd/bot/main.go # Entry point
├── internal/
│ ├── price/ # Price fetching logic (CoinGecko)
│ └── telegram/ # Telegram bot logic
│ ├── bot.go
│ ├── handlers.go
│ ├── ui.go
│ ├── state.go
│ └── coins.go
├── go.mod
├── go.sum
└── .env

```

---

## ⚙️ Setup

### 1. Clone the repository

```bash
git clone https://github.com/EduBarreira1212/currency-price-bot.git
cd currency-price-bot
```

### 2. Install dependencies

```bash
go mod tidy
```

### 3. Create a bot in Telegram

1. Open Telegram and search for [@BotFather](https://t.me/BotFather).
2. Send `/newbot` and follow the instructions.
3. Copy the **bot token**.

### 4. Set environment variables

Create a `.env` file:

```env
TELEGRAM_BOT_TOKEN=your_bot_token_here
```

### 5. Run the bot

```bash
go run cmd/bot/main.go
```

---

## 🛠 Usage

- `/start` → Shows the main menu
- Click **coin buttons** (💰 Get BTC, etc.) → Fetches the current price instantly
- Click **Start Updates** → Enables auto-updates
- Click **Stop Updates** → Disables auto-updates
- Click **interval buttons** → Changes how often prices are sent
- Click **currency buttons** → Switches between USD, EUR, BRL

---

## 📦 Deployment

You can deploy this bot on:

- [Railway](https://railway.app/) (easy & free tier)
- [Render](https://render.com/)
- [Fly.io](https://fly.io/)
- VPS (with `systemd` or Docker)

---

## ➕ Adding More Coins

Edit `internal/telegram/coins.go`:

```go
var Coins = []Coin{
    {ID: "bitcoin",  Label: "BTC", Emoji: "📈"},
    {ID: "ethereum", Label: "ETH", Emoji: "📉"},
    {ID: "solana",   Label: "SOL", Emoji: "⚡"},
    {ID: "cardano",  Label: "ADA", Emoji: "🔷"},
    {ID: "dogecoin", Label: "DOGE", Emoji: "🐶"},
}
```

Restart the bot — new coins will automatically appear in the menu and updates.

---

## ✅ To do

- Unit tests
- Deploy
- Customize bot in telegram

---

## 📝 License

MIT License © 2025 [Eduardo Barreira](https://github.com/EduBarreira1212)
