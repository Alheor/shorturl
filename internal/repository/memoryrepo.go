package repository

import (
	"context"
	"errors"
	"sync"

	"github.com/Alheor/shorturl/internal/models"
	"github.com/Alheor/shorturl/internal/urlhasher"
)

// MemoryRepo structure
type MemoryRepo struct {
	list map[string]string
	sync.RWMutex
}

// Add Добавить URL
func (fr *MemoryRepo) Add(ctx context.Context, user *models.User, name string) (string, error) {

	select {
	case <-ctx.Done():
		return ``, ctx.Err()
	default:
	}

	fr.Lock()
	defer fr.Unlock()

	//Обработка существующих URL
	for hash, el := range fr.list {
		if el == name {
			return ``, &models.UniqueErr{Err: errors.New("url already exists"), ShortKey: hash}
		}
	}

	hash := urlhasher.GetHash(name)
	fr.list[hash] = name

	return hash, nil
}

// AddBatch Добавить URL пачкой
func (fr *MemoryRepo) AddBatch(ctx context.Context, user *models.User, list *[]models.BatchEl) error {

	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}

	fr.Lock()
	defer fr.Unlock()

	for _, v := range *list {
		fr.list[v.ShortURL] = v.OriginalURL
	}

	return nil
}

// GetByShortName получить URL по короткому имени
func (fr *MemoryRepo) GetByShortName(ctx context.Context, user *models.User, name string) (string, error) {

	select {
	case <-ctx.Done():
		return ``, ctx.Err()
	default:
	}

	fr.RLock()
	defer fr.RUnlock()

	el, exists := fr.list[name]
	if !exists {
		return ``, nil
	}

	return el, nil
}

// IsReady готовность репозитория
func (fr *MemoryRepo) IsReady(ctx context.Context) bool {
	select {
	case <-ctx.Done():
		return false
	default:
	}

	return fr.list != nil
}

func (fr *MemoryRepo) RemoveByOriginalURL(ctx context.Context, user *models.User, url string) error {
	return errors.New(`method "Remove" from memory repository not supported`)
}
