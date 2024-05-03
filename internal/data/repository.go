package data

import (
	"github.com/benk-techworld/www-backend/internal/db"
	"github.com/benk-techworld/www-backend/internal/models"
)

type Repository struct {
	Article interface {
		Create(article *models.Article) error
	}
}

func NewRepo(db *db.DB) *Repository {
	return &Repository{
		Article: ArticleRepo{db},
	}
}
