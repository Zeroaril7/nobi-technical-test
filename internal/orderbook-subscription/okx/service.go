package okx

import (
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"strings"
)

func ProcessOrderBookData(okxMsg []byte, pair string) (OrderbookSubscriptionOKX, error) {

	pairSplit := strings.Split(pair, "/")
	if len(pairSplit) != 2 {
		log.Printf("Invalid pair format: %s", pair)
		return OrderbookSubscriptionOKX{}, fmt.Errorf("invalid pair format: %s", pair)
	}
	pair0, pair1 := pairSplit[0], pairSplit[1]

	if !isValidOrderBookMessage(okxMsg) {
		log.Printf("Ignored invalid OKX message for %s: %s", pair, string(okxMsg))
		return OrderbookSubscriptionOKX{}, nil
	}

	var okxSnapshot struct {
		Data []struct {
			Asks [][]string `json:"asks"`
			Bids [][]string `json:"bids"`
			Ts   string     `json:"ts"`
		} `json:"data"`
	}

	err := json.Unmarshal(okxMsg, &okxSnapshot)
	if err != nil {
		log.Printf("Error unmarshaling OKX data for %s: %v", pair, err)
		return OrderbookSubscriptionOKX{}, err
	}

	if len(okxSnapshot.Data) > 0 {
		data := okxSnapshot.Data[0]

		highestBid, lowestAsk := "0", "0"
		if len(data.Bids) > 0 {
			highestBid = data.Bids[0][0]
		}
		if len(data.Asks) > 0 {
			lowestAsk = data.Asks[0][0]
		}

		bidPrice, _ := strconv.ParseFloat(highestBid, 64)
		askPrice, _ := strconv.ParseFloat(lowestAsk, 64)
		midPrice := (bidPrice + askPrice) / 2

		okxData := NewOrderbookData(
			"okx",
			"Crypto:"+pair,
			pair0,
			pair1,
			lowestAsk,
			highestBid,
			fmt.Sprintf("%.8f", midPrice),
		)

		log.Printf("Successfully processed OKX order book data for %s", pair)
		return okxData, nil
	}

	log.Printf("No data available in OKX message for %s", pair)
	return OrderbookSubscriptionOKX{}, nil
}
