package urlhasher

import (
	"math/rand"
	"time"
)

// ShortNameLength string length
const shortNameLength = 8

var ShortNameGenerator RandomStringGenerator

type RandomStringGenerator interface {
	Generate() string
}

// ShortName short name structure
type ShortName struct{}

func Init() {
	ShortNameGenerator = new(ShortName)
}

// Generate создание хэша для URL
func (sh *ShortName) Generate() string {
	randSource := rand.New(rand.NewSource(time.Now().UnixNano()))
	charset := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

	b := make([]byte, shortNameLength)
	for i := range b {
		b[i] = charset[randSource.Intn(len(charset))]
	}

	return string(b)
}
