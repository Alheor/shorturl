package repository

import (
	"sync"

	"github.com/Alheor/shorturl/internal/urlhasher"
)

type UrlMap struct {
	list map[string]string
	sync.RWMutex
}

var urlMap *UrlMap

func GetRepository() *UrlMap {

	if urlMap == nil {
		urlMap = &UrlMap{list: make(map[string]string)}
	}

	return urlMap
}

// Add Добавить URL
func (sn *UrlMap) Add(URL string) string {

	sn.Lock()
	defer sn.Unlock()

	//Обработка существующих URL
	for hash, el := range urlMap.list {
		if el == URL {
			return hash
		}
	}

	//Уменьшить вероятность коллизии хэша
	hash := urlhasher.ShortNameGenerator.Generate()
	if _, exists := urlMap.list[hash]; exists {
		hash = urlhasher.ShortNameGenerator.Generate()
	}

	urlMap.list[hash] = URL

	return hash
}

// GetByShortName получить URL по короткому имени
func (sn *UrlMap) GetByShortName(name string) *string {

	sn.RLock()
	defer sn.RUnlock()

	URL, exists := urlMap.list[name]
	if !exists {
		return nil
	}

	return &URL
}
