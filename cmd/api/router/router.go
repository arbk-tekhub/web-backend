package router

import "github.com/gin-gonic/gin"

func Routes() *gin.Engine {

	router := gin.New()
	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	return router
}
