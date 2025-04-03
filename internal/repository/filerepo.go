package repository

import (
	"bufio"
	"context"
	"encoding/json"
	"os"
	"sync"

	"github.com/Alheor/shorturl/internal/logger"
	"github.com/Alheor/shorturl/internal/urlhasher"
)

var fileRepo *FileRepo

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
	for hash, el := range fileRepo.list {
		if el == name {
			return hash, nil
		}
	}

	//Уменьшить вероятность коллизии хэша
	hash := urlhasher.GetShortNameGenerator().Generate()
	if _, exists := fileRepo.list[hash]; exists {
		hash = urlhasher.GetShortNameGenerator().Generate()
	}

	fileRepo.list[hash] = name

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

// GetByShortName получить URL по короткому имени
func (fr *FileRepo) GetByShortName(ctx context.Context, name string) (string, error) {

	select {
	case <-ctx.Done():
		return ``, ctx.Err()
	default:
	}

	fr.RLock()
	defer fr.RUnlock()

	el, exists := fileRepo.list[name]
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

func load(ctx context.Context, path string) error {

	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}

	file, err := os.OpenFile(path, os.O_RDONLY|os.O_CREATE, 0666)
	if err != nil {
		return err
	}

	fileRepo.Lock()
	defer fileRepo.Unlock()

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		data := scanner.Bytes()

		el := URL{}
		err = json.Unmarshal(data, &el)
		if err != nil {
			continue
		}

		fileRepo.list[el.ID] = el.URL
	}

	err = file.Close()
	if err != nil {
		return err
	}

	return nil
}
