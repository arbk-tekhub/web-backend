package main

import (
	"net/http"
	"time"

	_ "github.com/benk-techworld/www-backend/docs"
	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func (app *application) Routes() *gin.Engine {

	router := gin.New()
	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	router.HandleMethodNotAllowed = true
	router.NoMethod(app.methodNotAllowedResponse)
	router.NoRoute(app.notFoundResponse)

	router.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"timestamp": time.Now().UnixNano(),
		})
	})

	// @Endpoint for LB/GW healthchecks
	router.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	// @Swagger private endpoint for API documentation
	router.GET("/docs/*any", app.requireBasicAuthentication(), ginSwagger.WrapHandler(swaggerfiles.Handler))

	// @Metadata endpoint
	router.GET("/latest/metadata", app.metaDataHandler)

	// @API V1:
	v1 := router.Group("/v1")
	{
		// @Public routes
		v1.GET("/articles/:id", app.fetchArticleHandler)
		v1.GET("/articles", app.fetchArticlesHandler)
		v1.POST("/subs", app.createSubscriberHandler)

		// @Private routes
		v1.Use(app.requireBasicAuthentication())

		v1.POST("/articles", app.createArticleHandler)
		v1.DELETE("/articles/:id", app.deleteArticleHandler)
		v1.PATCH("/articles/:id", app.updateArticleHandler)

		v1.GET("/subs", app.fetchSubsribersHandler)
		v1.DELETE("/subs/:id", app.deleteSubscriberHandler)
	}

	return router
}
