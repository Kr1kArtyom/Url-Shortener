package local

import (
	"context"
	"fmt"
	"sync"
	"url-shortener/internal/entity"
	"url-shortener/internal/repository"
)

type LocalRepository struct {
	originalURL map[string]struct{}
	hashTable   map[string]string
	mu          sync.RWMutex
}

func NewLocalRepository() *LocalRepository {
	return &LocalRepository{
		originalURL: make(map[string]struct{}),
		hashTable:   make(map[string]string),
	}
}

func (r *LocalRepository) InsertURL(ctx context.Context, u *entity.URL) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, ok := r.originalURL[u.OriginalURL]; ok {
		return fmt.Errorf("original url %s already exists", u)
	}

	r.hashTable[u.ShortURL] = u.OriginalURL
	r.originalURL[u.OriginalURL] = struct{}{}

	return nil
}

func (r *LocalRepository) GetOriginalURL(ctx context.Context, shortURL string) (string, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	originalURL, exists := r.hashTable[shortURL]
	if !exists {
		return "", repository.ErrURLNotFound
	}

	return originalURL, nil
}

func (r *LocalRepository) GetURLByOriginal(ctx context.Context, originalURL string) (*entity.URL, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	for shortURL, storedOriginalURL := range r.hashTable {
		if storedOriginalURL == originalURL {
			return &entity.URL{
				OriginalURL: originalURL,
				ShortURL:    shortURL,
			}, nil
		}
	}
	return nil, nil
}
