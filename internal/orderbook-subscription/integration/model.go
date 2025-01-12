package integration

type OrderbookSubscriptionData struct {
	Source      string `json:"source"`
	Pair        string `json:"pair"`
	Pair0       string `json:"pair0"`
	Pair1       string `json:"pair1"`
	LastUpdated int64  `json:"last_updated"`
	AskPrice    string `json:"ask_price"`
	BidPrice    string `json:"bid_price"`
	MidPrice    string `json:"mid_price"`
}

type RedisData struct {
	Binance OrderbookSubscriptionData `json:"binance"`
	OKX     OrderbookSubscriptionData `json:"okx"`
}

func NewRedisData(binance, okx OrderbookSubscriptionData) RedisData {
	return RedisData{
		Binance: binance,
		OKX:     okx,
	}
}
