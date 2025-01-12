package integration

import (
	"context"
	"log"
	"sync"
	"time"

	"github.com/Zeroaril7/nobi-technical-test/internal/orderbook-subscription/binance"
	"github.com/Zeroaril7/nobi-technical-test/internal/orderbook-subscription/crypto"
	"github.com/Zeroaril7/nobi-technical-test/internal/orderbook-subscription/okx"
	"github.com/gorilla/websocket"
)

var (
	reconnectChan = make(chan bool)
)

func fetchInitialPairs(service crypto.CryptoService, ctx context.Context) ([]string, error) {

	pairs, _, err := service.FindAll(ctx)
	if err != nil {
		return nil, err
	}

	var result []string
	for _, p := range pairs {
		result = append(result, p.Pair)
	}

	return result, nil
}

func startWebSocketsWithReconnect() {
	for {
		mu.Lock()
		currentPairs := activePairs
		mu.Unlock()

		log.Println("Starting WebSocket connections...")
		var wg sync.WaitGroup
		runningPairs := make(map[string]bool)

		for _, pair := range currentPairs {
			wg.Add(1)
			runningPairs[pair] = true
			go func(pair string) {
				defer wg.Done()
				startWebSocketForPair(pair)
			}(pair)
		}
		wg.Wait()

		log.Println("WebSocket connections established. Waiting for reconnect signal...")

		select {
		case <-reconnectChan:
			log.Println("Reconnect triggered. Restarting WebSocket connections...")
		case <-time.After(10 * time.Second):
			log.Println("No reconnect signal. Rechecking pairs...")
		}
	}
}

func startWebSocketForPair(pair string) {
	log.Printf("Starting WebSocket for pair: %s", pair)

	binanceConn, err := ConnectToBinanceWebSocket(formatPairForBinance(pair))
	if err != nil {
		log.Printf("Failed to connect to Binance WebSocket for %s: %v", pair, err)
		return
	}
	defer binanceConn.Close()

	okxConn, err := ConnectToOKXWebSocket(formatPairForOKX(pair))
	if err != nil {
		log.Printf("Failed to connect to OKX WebSocket for %s: %v", pair, err)
		return
	}
	defer okxConn.Close()

	listenToWebSocket(binanceConn, okxConn, pair)
}

func listenToWebSocket(binanceConn, okxConn *websocket.Conn, pair string) {
	for {
		select {
		case <-reconnectChan:
			log.Printf("Reconnect signal received. Closing WebSocket for %s", pair)
			return
		default:
			_, binanceMsg, err := binanceConn.ReadMessage()
			if err != nil {
				log.Printf("Error reading Binance WebSocket for %s: %v", pair, err)
				return
			}

			binanceData, err := binance.ProcessOrderBookData(binanceMsg, pair)
			if err != nil {
				log.Printf("Error processing Binance order book data for %s: %v", pair, err)
			}

			_, okxMsg, err := okxConn.ReadMessage()
			if err != nil {
				log.Printf("Error reading OKX WebSocket for %s: %v", pair, err)
				return
			}

			okxData, err := okx.ProcessOrderBookData(okxMsg, pair)
			if err != nil {
				log.Printf("Error processing OKX order book data for %s: %v", pair, err)
			}

			convertedOKXData := OrderbookSubscriptionData{
				Source:      okxData.Source,
				Pair:        okxData.Pair,
				Pair0:       okxData.Pair0,
				Pair1:       okxData.Pair1,
				LastUpdated: okxData.LastUpdated,
				AskPrice:    okxData.AskPrice,
				BidPrice:    okxData.BidPrice,
				MidPrice:    okxData.MidPrice,
			}

			convertedBinanceData := OrderbookSubscriptionData{
				Source:      binanceData.Source,
				Pair:        binanceData.Pair,
				Pair0:       binanceData.Pair0,
				Pair1:       binanceData.Pair1,
				LastUpdated: binanceData.LastUpdated,
				AskPrice:    binanceData.AskPrice,
				BidPrice:    binanceData.BidPrice,
				MidPrice:    binanceData.MidPrice,
			}

			SaveOrderbookToRedis(pair, convertedBinanceData, convertedOKXData)
		}
	}
}
