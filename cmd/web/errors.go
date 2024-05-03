package main

import (
	"fmt"
	"log/slog"
	"net/http"
	"runtime/debug"

	"github.com/gin-gonic/gin"
)

func (app *application) logError(r *http.Request, err error) {
	var (
		message string = err.Error()
		method  string = r.Method
		url     string = r.URL.String()
		trace   string = string(debug.Stack())
	)
	requestAttrs := slog.Group("request", "method", method, "url", url)
	app.logger.Error(message, requestAttrs, "trace", trace)
}

func (app *application) internalServerErrorResponse(c *gin.Context, err error) {
	message := "the server encountered a problem and could not process your request"
	app.logError(c.Request, err)
	c.JSON(http.StatusInternalServerError, gin.H{
		"error": message,
	})
}

func (app *application) notFoundResponse(c *gin.Context) {
	message := "the requested resource could not be found"
	c.JSON(http.StatusNotFound, gin.H{
		"error": message,
	})
}

func (app *application) methodNotAllowedResponse(c *gin.Context) {
	message := fmt.Sprintf("the %s method is not supported for this resource", c.Request.Method)
	c.JSON(http.StatusMethodNotAllowed, gin.H{
		"error": message,
	})
}

func (app *application) badRequestResponse(c *gin.Context, err error) {
	message := "The request could not be understood by the server due to malformed syntax or incorrect parameter type"
	app.logError(c.Request, err)
	c.JSON(http.StatusBadRequest, gin.H{
		"error": message,
	})
}
