package binance

import (
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"strings"
)

func ProcessOrderBookData(binanceMsg []byte, pair string) (binanceData OrderbookSubscriptionBinance, err error) {
	pairSplit := strings.Split(pair, "/")
	if len(pairSplit) != 2 {
		log.Printf("Invalid pair format: %s", pair)
		return OrderbookSubscriptionBinance{}, fmt.Errorf("invalid pair format")
	}
	pair0 := pairSplit[0]
	pair1 := pairSplit[1]

	var depthUpdate BinanceDepthUpdate
	err = json.Unmarshal(binanceMsg, &depthUpdate)
	if err != nil {
		log.Printf("Error unmarshaling Binance data for %s: %v", pair, err)
		return
	}

	if len(depthUpdate.Bids) == 0 || len(depthUpdate.Asks) == 0 {
		log.Printf("No bids or asks in Binance data for %s", pair)
		err = fmt.Errorf("no bids or asks in Binance data")
		return
	}

	highestBid := depthUpdate.Bids[0][0]
	lowestAsk := depthUpdate.Asks[0][0]

	bidPrice, _ := strconv.ParseFloat(highestBid, 64)
	askPrice, _ := strconv.ParseFloat(lowestAsk, 64)
	midPrice := (bidPrice + askPrice) / 2

	binanceData = NewOrderbookData(
		"binance",
		"Crypto:"+pair,
		pair0,
		pair1,
		lowestAsk,
		highestBid,
		fmt.Sprintf("%.8f", midPrice),
	)

	log.Printf("Successfully processed Binance order book data for %s", pair)
	return binanceData, nil
}
