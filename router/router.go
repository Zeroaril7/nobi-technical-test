package router

import (
	database "github.com/Zeroaril7/nobi-technical-test/pkg/databases"
	v1 "github.com/Zeroaril7/nobi-technical-test/router/v1"
	"github.com/gin-gonic/gin"
)

func SetupRoutes(app *gin.Engine) {
	db := database.DB.Instance

	v1Group := app.Group("/api/v1")
	v1.SetupV1Routes(v1Group, db)
}
