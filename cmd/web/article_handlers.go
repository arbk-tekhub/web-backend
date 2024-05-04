package main

import (
	"errors"
	"net/http"

	"github.com/benk-techworld/www-backend/internal/service"
	"github.com/benk-techworld/www-backend/internal/utils"
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

func (app *application) fetchArticleHandler(c *gin.Context) {

	idString := c.Param("id")

	article, err := app.service.GetArticleByID(idString)
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
		"article": article,
	})
}

func (app *application) deleteArticleHandler(c *gin.Context) {

	idString := c.Param("id")

	err := app.service.DeleteArticle(idString)
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
		"message": "article successfully deleted",
	})
}

func (app *application) fetchArticlesHandler(c *gin.Context) {

	var input service.FilterArticlesInput

	qs := c.Request.URL.Query()

	input.Title = utils.ReadStringFromQueryParams(qs, "title", "")
	input.Tags = utils.ReadCsvFromQueryParams(qs, "tags", []string{})
	input.Page = utils.ReadIntFromQueryParams(qs, "page", 1)
	input.PageSize = utils.ReadIntFromQueryParams(qs, "page_size", 20)
	input.Sort = utils.ReadStringFromQueryParams(qs, "sort", "-published")

	articles, err := app.service.GetArticles(&input)
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
		"articles": articles,
	})
}

func (app *application) updateArticleHandler(c *gin.Context) {

	idString := c.Param("id")

	var input service.UpdateArticleInput

	err := c.BindJSON(&input)
	if err != nil {
		app.badRequestResponse(c, err)
		return
	}

	updatedArticle, err := app.service.UpdateArticle(idString, &input)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrNotFound):
			app.notFoundResponse(c)
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
		"article": updatedArticle,
	})

}
