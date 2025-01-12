package v1

import (
	"github.com/Zeroaril7/nobi-technical-test/internal/orderbook-subscription/crypto"
	"github.com/Zeroaril7/nobi-technical-test/internal/orderbook-subscription/integration"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func SetupCryptoRoutes(router *gin.RouterGroup, db *gorm.DB) {

	repository := crypto.NewRepository(db)
	service := crypto.NewService(repository)
	handler := crypto.NewHandler(service)

	apiGroup := router.Group("/crypto")
	apiGroup.POST("", handler.Add)
	apiGroup.GET("", handler.FindAll)
	apiGroup.DELETE("/:id", handler.Delete)

	apiGroup.POST("/start-websocket", func(c *gin.Context) {
		integration.StartWebSocketHandler(c, service)
	})

	apiGroup.GET("/orderbook", func(c *gin.Context) {
		integration.FindByKey(c)
	})
}
