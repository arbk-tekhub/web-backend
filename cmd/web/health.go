package main

import (
	"net/http"
	"runtime"

	"github.com/benk-techworld/www-backend/internal/version"
	"github.com/gin-gonic/gin"
)

func (app *application) healthCheckHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":      "available",
		"environment": app.config.env,
		"build_info": map[string]string{
			"version":  version.Get(),
			"commitID": version.GetCommitID()[:7],
			"go":       runtime.Version(),
		},
	})
}
