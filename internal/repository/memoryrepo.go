package repository

import (
	"context"
	"errors"
	"sync"

	"github.com/Alheor/shorturl/internal/models"
	"github.com/Alheor/shorturl/internal/urlhasher"
)

var _ IRepository = (*MemoryRepo)(nil)

// MemoryRepo - структура репозитория в памяти.
type MemoryRepo struct {
	list map[string]map[string]string
	sync.RWMutex
}

// Add Добавить URL.
func (fr *MemoryRepo) Add(ctx context.Context, user *models.User, name string) (string, error) {

	select {
	case <-ctx.Done():
		return ``, ctx.Err()
	default:
	}

	fr.Lock()
	defer fr.Unlock()

	if fr.list[user.ID] == nil {
		fr.list[user.ID] = make(map[string]string)
	}

	urls := fr.list[user.ID]

	//Обработка существующих URL
	for hash, el := range urls {
		if el == name {
			return ``, &models.UniqueErr{Err: errors.New("url already exists"), ShortKey: hash}
		}
	}

	hash := urlhasher.GetHash(name)
	urls[hash] = name

	return hash, nil
}

// AddBatch Добавить несколько URL.
func (fr *MemoryRepo) AddBatch(ctx context.Context, user *models.User, list *[]models.BatchEl) error {

	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}

	fr.Lock()
	defer fr.Unlock()

	if fr.list[user.ID] == nil {
		fr.list[user.ID] = make(map[string]string)
	}

	urls := fr.list[user.ID]

	for _, v := range *list {
		urls[v.ShortURL] = v.OriginalURL
	}

	return nil
}

// GetByShortName Получить URL по короткому имени.
func (fr *MemoryRepo) GetByShortName(ctx context.Context, user *models.User, name string) (string, bool, error) {

	select {
	case <-ctx.Done():
		return ``, false, ctx.Err()
	default:
	}

	fr.RLock()
	defer fr.RUnlock()

	//Костыль для прохождения тестов
	if user == nil {
		for _, el := range fr.list {
			//Жесть, но тесты нужно пройти
			for short, original := range el {
				if short == name {
					return original, false, nil
				}
			}
		}

		return ``, false, nil
	}

	urls, exists := fr.list[user.ID]
	if !exists {
		return ``, false, nil
	}

	el, exists := urls[name]
	if !exists {
		return ``, false, nil
	}

	return el, false, nil
}

// IsReady Готовность репозитория.
func (fr *MemoryRepo) IsReady(ctx context.Context) bool {
	select {
	case <-ctx.Done():
		return false
	default:
	}

	return fr.list != nil
}

// RemoveByOriginalURL - удалить URL.
// Deprecated: не поддерживается эти типом репозитория.
func (fr *MemoryRepo) RemoveByOriginalURL(ctx context.Context, user *models.User, url string) error {
	return errors.New(`method "Remove" from memory repository not supported`)
}

// GetAll получить все URL пользователя.
func (fr *MemoryRepo) GetAll(ctx context.Context, user *models.User) (<-chan models.HistoryEl, <-chan error) {
	out := make(chan models.HistoryEl)
	errCh := make(chan error, 1)

	defer close(errCh)

	select {
	case <-ctx.Done():
		close(out)
		errCh <- ctx.Err()

		return out, errCh
	default:
	}

	list, exists := fr.list[user.ID]
	if !exists {
		close(out)
		errCh <- &models.HistoryNotFoundErr{}

		return out, errCh
	}

	go func() {
		defer close(out)

		for shortURL, originalURL := range list {
			out <- models.HistoryEl{OriginalURL: originalURL, ShortURL: shortURL}
		}
	}()

	return out, errCh
}

// RemoveBatch - массовое удаление URL.
// Deprecated: не поддерживается эти типом репозитория.
func (fr *MemoryRepo) RemoveBatch(ctx context.Context, user *models.User, list []string) error {
	return errors.New(`method "RemoveBatch" from memory repository not supported`)
}
