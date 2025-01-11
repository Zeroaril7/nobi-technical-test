package v1

import (
	"github.com/gin-gonic/gin"
)

func SetupV1Routes(router *gin.RouterGroup) {

	SetupEthereumFetcherRoutes(router)
}
