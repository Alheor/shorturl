package repository

import (
	"bufio"
	"encoding/json"
	"os"
	"sync"

	"github.com/Alheor/shorturl/internal/config"
	"github.com/Alheor/shorturl/internal/logger"
	"github.com/Alheor/shorturl/internal/urlhasher"
)

type URL struct {
	ID  string `json:"id"`
	URL string `json:"url"`
}

type URLMap struct {
	list map[string]string
	file *os.File
	sync.RWMutex
}

var urlMap *URLMap

func Init() error {
	if urlMap != nil {
		urlMap = nil
	}

	urlMap = &URLMap{list: make(map[string]string)}

	err := load(urlMap, config.GetOptions().FileStoragePath)
	if err != nil {
		return err
	}

	file, err := os.OpenFile(config.GetOptions().FileStoragePath, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		return err
	}

	urlMap.file = file

	return nil
}

func GetRepository() *URLMap {
	if urlMap == nil {
		logger.Fatal(`uninitialized repository`, nil)
	}

	return urlMap
}

// Add Добавить URL
func (sn *URLMap) Add(name string) (*string, error) {

	sn.Lock()
	defer sn.Unlock()

	//Обработка существующих URL
	for hash, el := range urlMap.list {
		if el == name {
			return &hash, nil
		}
	}

	//Уменьшить вероятность коллизии хэша
	hash := urlhasher.ShortNameGenerator.Generate()
	if _, exists := urlMap.list[hash]; exists {
		hash = urlhasher.ShortNameGenerator.Generate()
	}

	urlMap.list[hash] = name

	data, err := json.Marshal(&URL{ID: hash, URL: name})
	if err != nil {
		logger.Error(`marshal error`, err)
		return nil, err
	}

	data = append(data, '\n')

	_, err = sn.file.Write(data)
	if err != nil {
		logger.Error(`file write error`, err)
		return nil, err
	}

	return &hash, nil
}

// GetByShortName получить URL по короткому имени
func (sn *URLMap) GetByShortName(name string) *string {

	sn.RLock()
	defer sn.RUnlock()

	el, exists := urlMap.list[name]
	if !exists {
		return nil
	}

	return &el
}

// IsReady готовность репозитория
func (sn *URLMap) IsReady() bool {
	return sn.file != nil
}

func load(um *URLMap, path string) error {

	file, err := os.OpenFile(path, os.O_RDONLY|os.O_CREATE, 0666)
	if err != nil {
		return err
	}

	um.Lock()
	defer um.Unlock()

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		data := scanner.Bytes()

		el := URL{}
		err = json.Unmarshal(data, &el)
		if err != nil {
			continue
		}

		um.list[el.ID] = el.URL
	}

	err = file.Close()
	if err != nil {
		return err
	}

	return nil
}
