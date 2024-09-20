package repository

import (
	"context"
	"errors"
	"url-shortener/internal/entity"
)

type Repository interface {
	InsertURL(ctx context.Context, u *entity.URL) error
	GetOriginalURL(ctx context.Context, shortURL string) (string, error)
	GetURLByOriginal(ctx context.Context, originalURL string) (*entity.URL, error)
}

var (
	ErrURLNotFound = errors.New("url doesn't exist")
)
