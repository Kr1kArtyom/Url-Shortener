package services

import (
	"context"
	"errors"
	"log/slog"
	"url-shortener/internal/entity"
	"url-shortener/internal/repository"
	"url-shortener/internal/utils"
)

type URLShortener struct {
	logger *slog.Logger
	db     Database
}

func NewURLShortener(logger *slog.Logger, db Database) *URLShortener {
	return &URLShortener{logger, db}
}

func (u *URLShortener) CreateShortURL(ctx context.Context, originalURL string) (string, error) {
	u.logger.Info("Creating short URL", slog.String("originalURL", originalURL))

	existingURL, err := u.db.GetURLByOriginal(ctx, originalURL)
	if err == nil && existingURL != nil {
		u.logger.Info("Original URL already exists, returning existing short URL", slog.String("shortURL", existingURL.ShortURL))
		return existingURL.ShortURL, nil
	}

	increment := 0
	shortURL := utils.HashURL(originalURL, increment)

	for {
		_, err := u.db.GetOriginalURL(ctx, shortURL)
		if err != nil {
			if errors.Is(err, repository.ErrURLNotFound) {
				break
			}
			u.logger.Error("Error checking for URL uniqueness", slog.String("error", err.Error()))
			return "", err
		}

		u.logger.Info("Short URL collision detected, incrementing", slog.Int("increment", increment))
		increment++
		shortURL = utils.HashURL(originalURL, increment)
	}

	urlEntity := &entity.URL{
		OriginalURL: originalURL,
		ShortURL:    shortURL,
	}

	err = u.db.InsertURL(ctx, urlEntity)
	if err != nil {
		u.logger.Error("Failed to create short URL", slog.String("error", err.Error()))
		return "", err
	}

	return shortURL, nil
}

func (u *URLShortener) GetOriginalURL(ctx context.Context, shortURL string) (string, error) {
	u.logger.Info("Retrieving original URL", slog.String("shortURL", shortURL))

	originalURL, err := u.db.GetOriginalURL(ctx, shortURL)
	if err != nil {
		u.logger.Error("Failed to retrieve original URL", slog.String("error", err.Error()))
		return "", errors.New("URL not found")
	}

	return originalURL, nil
}
