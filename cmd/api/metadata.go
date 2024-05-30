package main

import (
	"net/http"
	"runtime"

	"github.com/benk-techworld/www-backend/internal/version"
	"github.com/gin-gonic/gin"
)

func (app *application) metaDataHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"environment": app.config.env,
		"system_info": map[string]string{
			"os":   runtime.GOOS,
			"arch": runtime.GOARCH,
		},
		"build_info": map[string]string{
			"version": version.Get(),
			"go":      runtime.Version(),
		},
	})
}
