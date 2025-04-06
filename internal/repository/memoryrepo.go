package repository

import (
	"context"
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
func (fr *MemoryRepo) Add(ctx context.Context, name string) (string, error) {

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
			return hash, nil
		}
	}

	//Уменьшить вероятность коллизии хэша
	hash := urlhasher.GetShortNameGenerator().Generate()
	if _, exists := fr.list[hash]; exists {
		hash = urlhasher.GetShortNameGenerator().Generate()
	}

	fr.list[hash] = name

	return hash, nil
}

// AddBatch Добавить URL пачкой
func (fr *MemoryRepo) AddBatch(ctx context.Context, list *[]models.BatchEl) error {

	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}

	fr.Lock()
	defer fr.Unlock()

	for _, v := range *list {
		//Уменьшить вероятность коллизии хэша
		if _, exists := fr.list[v.ShortURL]; exists {
			v.ShortURL = urlhasher.GetShortNameGenerator().Generate()
		}

		fr.list[v.ShortURL] = v.OriginalURL
	}

	return nil
}

// GetByShortName получить URL по короткому имени
func (fr *MemoryRepo) GetByShortName(ctx context.Context, name string) (string, error) {

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
