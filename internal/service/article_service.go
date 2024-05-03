package service

import (
	"strings"
	"time"

	"github.com/benk-techworld/www-backend/internal/models"
	"github.com/benk-techworld/www-backend/internal/validator"
)

type CreateArticleInput struct {
	Title            string            `json:"title"`
	Content          string            `json:"content"`
	Tags             []string          `json:"tags"`
	Author           string            `bson:"author"`
	Status           string            `json:"status"`
	ValidationErrors map[string]string `json:"-"`
}

func (svc Service) CreateArticle(input *CreateArticleInput) error {

	v := validator.New()
	input.ValidationErrors = v.Errors

	v.Check(input.Title != "", "title", "must be provided")
	v.Check(len(input.Title) >= 2, "title", "must be atleast 2 char long")
	v.Check(input.Author != "", "author", "must be provided")
	v.Check(input.Tags != nil, "tags", "must be provided")
	v.Check(len(input.Tags) > 0, "tags", "must not be empty")
	v.Check(validator.Unique(input.Tags), "tags", "must contain unique values")
	v.Check(input.Status != "", "status", "must be provided")
	v.Check(validator.PermittedValues(strings.ToLower(input.Status), "published", "unpublished"), "status", "unrecognized value")

	if v.HasErrors() {
		return ErrFailedValidation
	}

	article := models.Article{
		Title:   input.Title,
		Content: input.Content,
		Tags:    input.Tags,
		Author:  input.Author,
		Status:  strings.ToLower(input.Status),
		Created: time.Now(),
	}

	if strings.ToLower(article.Status) == "published" {
		article.Published = time.Now()
	}

	return svc.Repo.Article.Create(&article)
}
