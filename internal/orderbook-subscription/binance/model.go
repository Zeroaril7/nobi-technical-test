package binance

import "time"

type OrderbookSubscriptionBinance struct {
	Source      string `json:"source"`
	Pair        string `json:"pair"`
	Pair0       string `json:"pair0"`
	Pair1       string `json:"pair1"`
	LastUpdated int64  `json:"last_updated"`
	AskPrice    string `json:"ask_price"`
	BidPrice    string `json:"bid_price"`
	MidPrice    string `json:"mid_price"`
}

func NewOrderbookData(source, pair, pair0, pair1, askPrice, bidPrice, midPrice string) OrderbookSubscriptionBinance {
	return OrderbookSubscriptionBinance{
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

type BinanceDepthUpdate struct {
	Event          string     `json:"e"`
	EventTime      int64      `json:"E"`
	Symbol         string     `json:"s"`
	FirstUpdateID  int64      `json:"U"`
	FinalUpdateID  int64      `json:"u"`
	PreviousUpdate int64      `json:"pu"`
	Bids           [][]string `json:"b"`
	Asks           [][]string `json:"a"`
}
