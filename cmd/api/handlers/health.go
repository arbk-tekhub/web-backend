package handlers

import (
	"net/http"
	"runtime"

	"github.com/benk-techworld/www-backend/internal/version"
	"github.com/gin-gonic/gin"
)

func HealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status": "available",
		"build_info": map[string]string{
			"version": version.Get(),
			"commit":  version.GetCommitID(),
			"go":      runtime.Version(),
		},
	})
}
