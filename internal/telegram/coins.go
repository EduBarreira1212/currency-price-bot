package telegram

type Coin struct {
	ID    string
	Label string
	Emoji string
}

var Coins = []Coin{
	{ID: "bitcoin", Label: "BTC", Emoji: "📈"},
	{ID: "ethereum", Label: "ETH", Emoji: "📉"},
	{ID: "solana", Label: "SOL", Emoji: "⚡"},
	{ID: "cardano", Label: "ADA", Emoji: "🔷"},
	{ID: "dogecoin", Label: "DOGE", Emoji: "🐶"},
}

func FindCoin(id string) (Coin, bool) {
	for _, c := range Coins {
		if c.ID == id {
			return c, true
		}
	}
	return Coin{}, false
}
