package ethereumfetcher

import (
	"encoding/json"
	"net/http"

	redissdk "github.com/Zeroaril7/nobi-technical-test/pkg/redis-sdk"
	"github.com/gin-gonic/gin"
)

func GetExchangeRate(c *gin.Context) {
	key := "price:Crypto:APEETH/ETH"

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

	c.JSON(http.StatusOK, gin.H{"data": parsedData, "message": "Success get exchange rate"})
}
