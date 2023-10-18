// Package randomname
// Random string generator
package randomname

import (
	"math/rand"
	"time"
)

// RandomStringGenerator interface
type RandomStringGenerator interface {
	Generate() string
}

const shortNameLength = 8

// ShortName short name structure
type ShortName struct{}

// Init randomname constructor
func Init() RandomStringGenerator {
	return new(ShortName)
}

func (rg ShortName) Generate() string {
	randSource := rand.New(rand.NewSource(time.Now().UnixNano()))
	charset := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

	b := make([]byte, shortNameLength)
	for i := range b {
		b[i] = charset[randSource.Intn(len(charset))]
	}

	return string(b)
}
