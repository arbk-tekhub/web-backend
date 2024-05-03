package service

import "github.com/benk-techworld/www-backend/internal/data"

type Service struct {
	Repo *data.Repository
}

func New(r *data.Repository) *Service {
	return &Service{
		Repo: r,
	}
}
