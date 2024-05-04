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
	Author           string            `bson:"author"`
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

	if v.HasErrors() {
		return ErrFailedValidation
	}

	article := &models.Article{
		Title:   input.Title,
		Content: input.Content,
		Tags:    input.Tags,
		Author:  input.Author,
		Created: time.Now(),
	}

	// Default values
	article.Status = "published"

	return svc.Repo.Article.Create(article)
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

type FilterArticlesInput struct {
	Title string
	Tags  []string
	models.Filters
	ValidationErrors map[string]string
}

func (svc *Service) GetArticles(inputFilters *FilterArticlesInput) ([]models.Article, error) {

	v := validator.New()
	inputFilters.ValidationErrors = v.Errors
	inputFilters.SortSafeList = []string{"title", "published", "-title", "-published"}
	if models.ValidateFilters(v, inputFilters.Filters); v.HasErrors() {
		return nil, ErrFailedValidation
	}

	articles, err := svc.Repo.Article.Get(inputFilters.Title, inputFilters.Tags, inputFilters.Filters)
	if err != nil {
		return nil, err
	}

	return articles, nil
}

type UpdateArticleInput struct {
	Title            string            `bson:"title,omitempty" json:"title,omitempty"`
	Content          string            `bson:"content,omitempty" json:"content,omitempty"`
	Tags             []string          `bson:"tags,omitempty" json:"tags,omitempty"`
	Author           string            `bson:"author,omitempty" json:"author,omitempty"`
	Status           string            `bson:"status,omitempty" json:"status,omitempty"`
	ValidationErrors map[string]string `bson:"-" json:"-"`
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

	v.Check(validator.PermittedValues(strings.ToLower(input.Status), "published", "unpublished"), "status", "unrecognized value")

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
