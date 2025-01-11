package router

import (
	v1 "github.com/Zeroaril7/nobi-technical-test/router/v1"
	"github.com/gin-gonic/gin"
)

func SetupRoutes(app *gin.Engine) {

	v1Group := app.Group("/api/v1")
	v1.SetupV1Routes(v1Group)
}
