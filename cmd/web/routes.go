package main

import (
	"net/http"

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

	// @Server side rendering
	router.Static("/static", "./assets/static/")
	router.LoadHTMLGlob("assets/templates/*")
	router.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", gin.H{
			"msg": "Benk Techworld Backend",
		})
	})

	// @Endpoint for load balancers
	router.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	// @Swagger private endpoint for API documentation
	router.GET("/docs/*any", app.requireBasicAuthentication(), ginSwagger.WrapHandler(swaggerfiles.Handler))

	// @API V1:
	v1 := router.Group("/v1")
	{
		// @Public routes
		v1.GET("/health", app.healthCheckHandler)
		v1.GET("/articles/:id", app.fetchArticleHandler)
		v1.GET("/articles", app.fetchArticlesHandler)

		// @Private routes
		v1.Use(app.requireBasicAuthentication())

		v1.POST("/articles", app.createArticleHandler)
		v1.DELETE("/articles/:id", app.deleteArticleHandler)
		v1.PATCH("/articles/:id", app.updateArticleHandler)
	}

	return router
}
