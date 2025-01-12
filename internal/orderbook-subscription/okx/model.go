package okx

import "time"

type OrderbookSubscriptionOKX struct {
	Source      string `json:"source"`
	Pair        string `json:"pair"`
	Pair0       string `json:"pair0"`
	Pair1       string `json:"pair1"`
	LastUpdated int64  `json:"last_updated"`
	AskPrice    string `json:"ask_price"`
	BidPrice    string `json:"bid_price"`
	MidPrice    string `json:"mid_price"`
}

func NewOrderbookData(source, pair, pair0, pair1, askPrice, bidPrice, midPrice string) OrderbookSubscriptionOKX {
	return OrderbookSubscriptionOKX{
		Source:      source,
		Pair:        pair,
		Pair0:       pair0,
		Pair1:       pair1,
		LastUpdated: time.Now().Unix(),
		AskPrice:    askPrice,
		BidPrice:    bidPrice,
		MidPrice:    midPrice,
	}
}

type OrderBookData struct {
	BidPrice string `json:"bidPrice"`
	AskPrice string `json:"askPrice"`
}
