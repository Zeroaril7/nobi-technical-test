package integration

import (
	"encoding/json"
	"log"
	"time"

	"github.com/Zeroaril7/nobi-technical-test/config"
	"github.com/Zeroaril7/nobi-technical-test/pkg/constant"
	redissdk "github.com/Zeroaril7/nobi-technical-test/pkg/redis-sdk"
	"github.com/gorilla/websocket"
)

func ConnectToBinanceWebSocket(symbol string) (*websocket.Conn, error) {
	url := constant.BINANCE_WEBSOCKET + symbol + "@depth5"
	conn, _, err := websocket.DefaultDialer.Dial(url, nil)
	if err != nil {
		log.Printf("Error connecting to Binance WebSocket: %v", err)
		return nil, err
	}

	go StartPingPong(conn, 25*time.Second, "binance")

	log.Printf("Successfully connected to binance WebSocket for symbol: %s", symbol)
	return conn, nil
}

func ConnectToOKXWebSocket(symbol string) (*websocket.Conn, error) {
	url := constant.OKX_WEBSOCKET + config.Config("OKX_API_KEY")
	conn, _, err := websocket.DefaultDialer.Dial(url, nil)
	if err != nil {
		log.Printf("Error connecting to OKX WebSocket: %v", err)
		return nil, err
	}

	subscribeMessage := map[string]interface{}{
		"op": "subscribe",
		"args": []map[string]string{
			{
				"channel": "books5",
				"instId":  symbol,
			},
		},
	}
	err = conn.WriteJSON(subscribeMessage)
	if err != nil {
		log.Printf("Error sending subscription message to OKX WebSocket: %v", err)
		return nil, err
	}

	go StartPingPong(conn, 25*time.Second, "okx")

	log.Printf("Successfully connected to OKX WebSocket for symbol: %s", symbol)
	return conn, nil
}

func SaveOrderbookToRedis(pair string, binance, okx OrderbookSubscriptionData) error {
	redisData := NewRedisData(binance, okx)

	dataJSON, err := json.Marshal(redisData)
	if err != nil {
		return err
	}

	key := "price:Crypto:" + pair
	err = redissdk.RedisClient.Set(redissdk.Ctx, key, dataJSON, 5*time.Minute).Err()
	if err != nil {
		return err
	}

	return nil
}

func StartPingPong(conn *websocket.Conn, interval time.Duration, websocketName string) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for range ticker.C {
		err := conn.WriteMessage(websocket.TextMessage, []byte("ping"))
		if err != nil {
			log.Printf("Error sending ping: %v", err)
			return
		}
		log.Println("Ping sent to ", websocketName, "Websocket")
	}
}
