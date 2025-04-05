package urlhasher

import (
	"math/rand"
	"time"
)

// ShortNameLength string length
const ShortNameLength = 8

var generator RandomStringGenerator

type RandomStringGenerator interface {
	Generate() string
}

// ShortName short name structure
type ShortName struct{}

func Init(g RandomStringGenerator) {

	if g != nil {
		generator = g
		return
	}

	generator = new(ShortName)
}

// Generate создание хэша для URL
func (sh *ShortName) Generate() string {
	randSource := rand.New(rand.NewSource(time.Now().UnixNano()))
	charset := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

	b := make([]byte, ShortNameLength)
	for i := range b {
		b[i] = charset[randSource.Intn(len(charset))]
	}

	return string(b)
}

func GetShortNameGenerator() RandomStringGenerator {
	return generator
}
