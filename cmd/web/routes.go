package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (app *application) Routes() *gin.Engine {

	router := gin.New()
	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	router.HandleMethodNotAllowed = true
	router.NoMethod(app.methodNotAllowedResponse)
	router.NoRoute(app.notFoundResponse)

	router.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	v1 := router.Group("/v1")
	v1.GET("/health", app.healthCheckHandler)

	v1.POST("/articles", app.createArticleHandler)

	return router
}
