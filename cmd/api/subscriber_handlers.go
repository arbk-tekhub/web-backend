package main

import (
	"errors"
	"net/http"

	"github.com/benk-techworld/www-backend/internal/service"
	"github.com/benk-techworld/www-backend/internal/utils"
	"github.com/gin-gonic/gin"
)

func (app *application) createSubscriberHandler(c *gin.Context) {

	var input service.CreateSubscriberInput

	err := c.BindJSON(&input)
	if err != nil {
		app.badRequestResponse(c, err)
		return
	}

	sub, err := app.service.CreateSubscriber(&input)
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
		"subscriber": sub,
	})
}

func (app *application) fetchSubsribersHandler(c *gin.Context) {

	var input service.FetchSubscribersInput

	qs := c.Request.URL.Query()

	input.Email = utils.ReadStringFromQueryParams(qs, "email", "")
	input.Page = utils.ReadIntFromQueryParams(qs, "page", 1)
	input.PageSize = utils.ReadIntFromQueryParams(qs, "page_size", 20)
	input.Sort = utils.ReadStringFromQueryParams(qs, "sort", "-created")

	subs, err := app.service.GetSubscribers(&input)
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

	c.JSON(http.StatusOK, gin.H{
		"subscribers": subs,
	})
}

func (app *application) deleteSubscriberHandler(c *gin.Context) {

	idString := c.Param("id")
	err := app.service.DeleteSubscriber(idString)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrNotFound):
			app.notFoundResponse(c)
		default:
			app.internalServerErrorResponse(c, err)
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "subscriber successfully deleted",
	})
}
