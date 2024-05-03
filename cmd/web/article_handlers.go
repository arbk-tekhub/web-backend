package main

import (
	"errors"
	"net/http"

	"github.com/benk-techworld/www-backend/internal/service"
	"github.com/gin-gonic/gin"
)

func (app *application) createArticleHandler(c *gin.Context) {

	var input service.CreateArticleInput

	err := c.BindJSON(&input)
	if err != nil {
		app.badRequestResponse(c, err)
		return
	}

	err = app.service.CreateArticle(&input)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrFailedValidation):
			c.JSON(http.StatusUnprocessableEntity, gin.H{
				"error": input.ValidationErrors,
			})
		default:
			app.internalServerErrorResponse(c, err)
		}
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "article successfully created",
	})

}
