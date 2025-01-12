package v1

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func SetupV1Routes(router *gin.RouterGroup, db *gorm.DB) {
	SetupEthereumFetcherRoutes(router)
	SetupCryptoRoutes(router, db)
}
