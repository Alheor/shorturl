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
	ID  string `json:"id"`
	URL string `json:"url"`
}

// FileRepo structure
type FileRepo struct {
	list map[string]string
	file *os.File
	sync.RWMutex
}

// Add Добавить URL
func (fr *FileRepo) Add(ctx context.Context, name string) (string, error) {

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

	//Уменьшить вероятность коллизии хэша
	hash := urlhasher.GetShortNameGenerator().Generate()
	if _, exists := fr.list[hash]; exists {
		hash = urlhasher.GetShortNameGenerator().Generate()
	}

	fr.list[hash] = name

	data, err := json.Marshal(&URL{ID: hash, URL: name})
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
func (fr *FileRepo) AddBatch(ctx context.Context, list *[]models.BatchEl) error {

	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}

	fr.Lock()
	defer fr.Unlock()

	var data []byte
	var err error
	snGen := urlhasher.GetShortNameGenerator()

	for _, v := range *list {
		//Уменьшить вероятность коллизии хэша
		if _, exists := fr.list[v.ShortURL]; exists {
			v.ShortURL = snGen.Generate()
		}

		el, err := json.Marshal(&URL{ID: v.ShortURL, URL: v.OriginalURL})
		if err != nil {
			return err
		}

		fr.list[v.ShortURL] = v.OriginalURL
		data = append(data, append(el, '\n')...)
	}

	_, err = fr.file.Write(data)
	if err != nil {
		return err
	}

	return nil
}

// GetByShortName получить URL по короткому имени
func (fr *FileRepo) GetByShortName(ctx context.Context, name string) (string, error) {

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
func (fr *FileRepo) IsReady(ctx context.Context) bool {
	select {
	case <-ctx.Done():
		return false
	default:
	}

	return fr.file != nil
}

func (fr *FileRepo) RemoveByOriginalUrl(ctx context.Context, url string) error {
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

		fr.list[el.ID] = el.URL
	}

	return nil
}
