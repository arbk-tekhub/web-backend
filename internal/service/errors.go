package service

import "errors"

var (
	ErrFailedValidation = errors.New("failed validation")
	ErrNotFound         = errors.New("document not found")
)
