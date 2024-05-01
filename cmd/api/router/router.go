package router

import (
	"github.com/benk-techworld/www-backend/cmd/api/handlers"
	"github.com/gin-gonic/gin"
)

func Routes() *gin.Engine {

	router := gin.New()
	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	v1 := router.Group("/v1")
	v1.GET("/health", handlers.HealthCheck)

	return router
}
