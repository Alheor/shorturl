// Package urlhasher - сервис хеширования URL.
//
// # Описание
//
// Хеширует оригинальный URL и возвращает его сокращенное представление.
package urlhasher

import (
	"strconv"

	"github.com/spaolacci/murmur3"
)

// HashLength - максимальная длинна хэша.
const HashLength = 20

// GetHash Получение сокращенного варианта URL.
func GetHash(URL string) string {
	m := murmur3.Sum64([]byte(URL))
	return strconv.FormatUint(m, 10)
}
