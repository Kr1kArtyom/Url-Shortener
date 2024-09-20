package postgres

import (
	"context"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"url-shortener/internal/entity"
	"url-shortener/internal/repository"
)

type PgRepository struct {
	conn *pgxpool.Pool
}

func NewPgRepository(conn *pgxpool.Pool) *PgRepository {
	return &PgRepository{
		conn: conn,
	}
}

func (r *PgRepository) InsertURL(ctx context.Context, u *entity.URL) error {
	_, err := r.conn.Exec(ctx,
		"INSERT INTO short_urls (short_url, long_url) VALUES($1, $2)",
		u.ShortURL, u.OriginalURL)
	if err != nil {
		return fmt.Errorf("unable to insert short URL: %v", err)
	}

	return nil
}

func (r *PgRepository) GetOriginalURL(ctx context.Context, shortURL string) (string, error) {
	var originalURL string
	err := r.conn.QueryRow(ctx,
		"SELECT long_url FROM short_urls WHERE short_url = $1",
		shortURL).Scan(&originalURL)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return "", repository.ErrURLNotFound
		}
		return "", err
	}

	return originalURL, nil
}

func (r *PgRepository) GetURLByOriginal(ctx context.Context, originalURL string) (*entity.URL, error) {
	row := r.conn.QueryRow(ctx, "SELECT short_url FROM short_urls WHERE long_url=$1", originalURL)
	url := &entity.URL{}
	if err := row.Scan(&url.ShortURL); err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil // URL not found
		}
		return nil, err
	}
	return url, nil
}
