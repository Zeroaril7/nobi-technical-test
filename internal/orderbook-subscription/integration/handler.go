package integration

import (
	"encoding/json"
	"log"
	"net/http"
	"sync"

	"github.com/Zeroaril7/nobi-technical-test/internal/orderbook-subscription/crypto"
	redissdk "github.com/Zeroaril7/nobi-technical-test/pkg/redis-sdk"
	"github.com/gin-gonic/gin"
)

var (
	activePairs []string
	mu          sync.Mutex
)

func StartWebSocketHandler(c *gin.Context, service crypto.CryptoService) {

	c.JSON(200, gin.H{"message": "WebSocket connections started"})

	pairs, err := fetchInitialPairs(service, c)
	if err != nil {
		log.Printf("Failed to fetch initial pairs: %v", err)
		return
	}

	mu.Lock()
	activePairs = pairs
	mu.Unlock()

	go startWebSocketsWithReconnect()

	go listenForNewPairs(service, c)
}

func FindByKey(c *gin.Context) {
	key := c.Query("key")

	data, err := redissdk.RedisClient.Get(redissdk.Ctx, key).Result()
	if err != nil {
		if err.Error() == "redis: nil" {
			c.JSON(http.StatusNotFound, gin.H{"error": "Data not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch data from Redis"})
		return
	}

	var parsedData map[string]interface{}
	err = json.Unmarshal([]byte(data), &parsedData)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse data"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": parsedData, "message": "Success get orderbook"})
}
