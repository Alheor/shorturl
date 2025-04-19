package repository

import (
	"bufio"
	"context"
	"encoding/json"
	"errors"
	"os"
	"sync"

	"github.com/Alheor/shorturl/internal/logger"
	"github.com/Alheor/shorturl/internal/models"
	"github.com/Alheor/shorturl/internal/urlhasher"
)

type URL struct {
	UserID string `json:"user_id"`
	ID     string `json:"id"`
	URL    string `json:"url"`
}

// FileRepo structure
type FileRepo struct {
	list map[string]map[string]string
	file *os.File
	sync.RWMutex
}

// Add Добавить URL
func (fr *FileRepo) Add(ctx context.Context, user *models.User, name string) (string, error) {

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

	data, err := json.Marshal(&URL{UserID: user.ID, ID: hash, URL: name})
	if err != nil {
		logger.Error(`marshal error`, err)
		return ``, err
	}

	data = append(data, '\n')

	_, err = fr.file.Write(data)
	if err != nil {
		logger.Error(`file write error`, err)
		return ``, err
	}

	return hash, nil
}

// AddBatch Добавить URL пачкой
func (fr *FileRepo) AddBatch(ctx context.Context, user *models.User, list *[]models.BatchEl) error {

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

	var data []byte
	var err error

	for _, v := range *list {
		v.ShortURL = urlhasher.GetHash(v.OriginalURL)

		el, err := json.Marshal(&URL{UserID: user.ID, ID: v.ShortURL, URL: v.OriginalURL})
		if err != nil {
			return err
		}

		urls[v.ShortURL] = v.OriginalURL
		data = append(data, append(el, '\n')...)
	}

	_, err = fr.file.Write(data)
	if err != nil {
		return err
	}

	return nil
}

// GetByShortName получить URL по короткому имени
func (fr *FileRepo) GetByShortName(ctx context.Context, user *models.User, name string) (string, error) {

	select {
	case <-ctx.Done():
		return ``, ctx.Err()
	default:
	}

	fr.RLock()
	defer fr.RUnlock()

	urls, exists := fr.list[user.ID]
	if !exists {
		return ``, nil
	}

	el, exists := urls[name]
	if !exists {
		return ``, nil
	}

	return el, nil
}

// IsReady готовность репозитория
func (fr *FileRepo) IsReady(ctx context.Context) bool {
	select {
	case <-ctx.Done():
		return false
	default:
	}

	return fr.file != nil
}

func (fr *FileRepo) RemoveByOriginalURL(ctx context.Context, user *models.User, url string) error {
	return errors.New(`method "Remove" from file repository not supported`)
}

func (fr *FileRepo) Load(ctx context.Context, path string) error {

	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}

	var err error

	fr.file, err = os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		return err
	}

	fr.Lock()
	defer fr.Unlock()

	scanner := bufio.NewScanner(fr.file)

	for scanner.Scan() {
		data := scanner.Bytes()

		el := URL{}
		err = json.Unmarshal(data, &el)
		if err != nil {
			continue
		}

		if fr.list[el.UserID] == nil {
			fr.list[el.UserID] = make(map[string]string)
		}

		fr.list[el.UserID][el.ID] = el.URL
	}

	return nil
}
