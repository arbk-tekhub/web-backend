package service

import (
	"errors"
	"strings"
	"time"

	"github.com/benk-techworld/www-backend/internal/models"
	"github.com/benk-techworld/www-backend/internal/utils"
	"github.com/benk-techworld/www-backend/internal/validator"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type CreateArticleInput struct {
	Title            string            `json:"title"`
	Content          string            `json:"content"`
	Tags             []string          `json:"tags"`
	Author           string            `json:"author"`
	ValidationErrors map[string]string `json:"-"`
}

type UpdateArticleInput struct {
	Title            string            `bson:"title,omitempty" json:"title,omitempty"`
	Content          string            `bson:"content,omitempty" json:"content,omitempty"`
	Tags             []string          `bson:"tags,omitempty" json:"tags,omitempty"`
	Author           string            `bson:"author,omitempty" json:"author,omitempty"`
	Status           string            `bson:"status,omitempty" json:"status,omitempty"`
	ValidationErrors map[string]string `bson:"-" json:"-"`
}

type FetchArticlesInput struct {
	Title string
	Tags  []string
	models.Filters
	ValidationErrors map[string]string
}

func (svc Service) CreateArticle(input *CreateArticleInput) (*models.Article, error) {

	v := validator.New()
	input.ValidationErrors = v.Errors

	v.Check(validator.NotBlank(input.Title), "title", "must be provided")
	v.Check(len(input.Title) >= 2, "title", "must be atleast 2 char long")
	v.Check(validator.NotBlank(input.Author), "author", "must be provided")
	v.Check(input.Tags != nil, "tags", "must be provided")
	v.Check(len(input.Tags) > 0, "tags", "must not be empty")
	v.Check(validator.Unique(input.Tags), "tags", "must contain unique values")

	if v.HasErrors() {
		return nil, ErrFailedValidation
	}

	defaultStatus := "published"

	article := &models.Article{
		Title:   input.Title,
		Content: input.Content,
		Tags:    input.Tags,
		Author:  input.Author,
		Status:  defaultStatus,
		Created: time.Now(),
	}

	err := svc.Repo.Article.Create(article)
	if err != nil {
		return nil, err
	}
	return article, nil
}

func (svc *Service) GetArticleByID(id string) (*models.Article, error) {

	docID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		if errors.Is(err, primitive.ErrInvalidHex) {
			return nil, ErrNotFound
		}
		return nil, err
	}

	article, err := svc.Repo.Article.GetByID(docID)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, ErrNotFound
		}
		return nil, err

	}
	return article, nil
}

func (svc *Service) DeleteArticle(id string) error {

	docID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		if errors.Is(err, primitive.ErrInvalidHex) {
			return ErrNotFound
		}
		return err
	}

	res, err := svc.Repo.Article.Delete(docID)
	if err != nil {
		return err
	}

	if res.DeletedCount == 0 {
		return ErrNotFound
	}

	return nil
}

func (svc *Service) GetArticles(input *FetchArticlesInput) ([]models.Article, error) {

	v := validator.New()
	input.ValidationErrors = v.Errors

	input.SortSafeList = []string{"title", "published", "-title", "-published"}
	if models.ValidateFilters(v, input.Filters); v.HasErrors() {
		return nil, ErrFailedValidation
	}

	articles, err := svc.Repo.Article.Get(input.Title, input.Tags, input.Filters)
	if err != nil {
		return nil, err
	}

	return articles, nil
}

func (svc *Service) UpdateArticle(id string, input *UpdateArticleInput) (*models.Article, error) {

	docID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		if errors.Is(err, primitive.ErrInvalidHex) {
			return nil, ErrNotFound
		}
		return nil, err
	}

	v := validator.New()
	input.ValidationErrors = v.Errors

	if validator.NotBlank(input.Title) {
		v.Check(len(input.Title) >= 2, "title", "must be atleast 2 char long")
	}

	if input.Tags != nil {
		v.Check(len(input.Tags) > 0, "tags", "must not be empty")
		v.Check(validator.Unique(input.Tags), "tags", "must contain unique values")
	}

	if validator.NotBlank(input.Title) {
		v.Check(validator.PermittedValues(strings.ToLower(input.Status), "published", "unpublished"), "status", "unrecognized value")
	}

	if v.HasErrors() {
		return nil, ErrFailedValidation
	}

	updateDoc, err := utils.StructToBsonMap(input)
	if err != nil {
		return nil, err
	}

	res := svc.Repo.Article.Update(docID, updateDoc)
	if res.Err() != nil {
		if errors.Is(res.Err(), mongo.ErrNoDocuments) {
			return nil, ErrNotFound
		}
		return nil, res.Err()
	}

	var updatedArticle models.Article

	if err = res.Decode(&updatedArticle); err != nil {
		return nil, err
	}

	return &updatedArticle, nil

}
