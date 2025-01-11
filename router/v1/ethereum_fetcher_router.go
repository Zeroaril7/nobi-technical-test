package v1

import (
	ethereumfetcher "github.com/Zeroaril7/nobi-technical-test/internal/ethereum-fetcher"
	"github.com/gin-gonic/gin"
)

func SetupEthereumFetcherRoutes(router *gin.RouterGroup) {

	apiGroup := router.Group("/ethereum")
	apiGroup.GET("/fetch-exchange-rate", ethereumfetcher.GetExchangeRate)

}
