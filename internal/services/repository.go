package services

import (
	"context"
	"url-shortener/internal/entity"
)

type Database interface {
	InsertURL(ctx context.Context, u *entity.URL) error
	GetOriginalURL(ctx context.Context, shortURL string) (string, error)
	GetURLByOriginal(ctx context.Context, originalURL string) (*entity.URL, error)
}
