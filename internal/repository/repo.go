package repository

import (
	"github.com/Alheor/shorturl/internal/urlhasher"
)

var urlMap map[string]string

func Init() {
	urlMap = make(map[string]string)
}

// Add Добавить URL
func Add(URL string) string {

	//Обработка существующих URL
	for hash, el := range urlMap {
		if el == URL {
			return hash
		}
	}

	//Уменьшит вероятность коллизии хэша
	hash := urlhasher.ShortNameGenerator.Generate()
	if _, exists := urlMap[hash]; !exists {
		hash = urlhasher.ShortNameGenerator.Generate()
	}

	urlMap[hash] = URL

	return hash
}

// GetByShortName получить URL по короткому имени
func GetByShortName(name string) *string {
	URL, exists := urlMap[name]
	if !exists {
		return nil
	}

	return &URL
}
