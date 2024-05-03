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

	router.Static("/static", "./assets/static/")

	router.LoadHTMLGlob("assets/templates/*")

	router.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", gin.H{
			"msg": "Benk Techworld Backend",
		})
	})

	router.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	v1 := router.Group("/v1")
	// @API V1: Public routes
	v1.GET("/health", app.healthCheckHandler)
	v1.GET("/articles/:id", app.fetchArticleHandler)

	// @API V1: Private routes
	v1.POST("/articles", app.createArticleHandler)
	v1.DELETE("/articles/:id", app.deleteArticleHandler)

	return router
}
