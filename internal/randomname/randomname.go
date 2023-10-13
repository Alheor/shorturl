// Package randomname
// Random string generator
package randomname

import (
	"math/rand"
	"time"
)

const shortNameLength = 8

// RandomString interface
type RandomString interface {
	Generate() string
}

// ShortName short name structure
type ShortName struct{}

func (rg ShortName) Generate() string {
	randSource := rand.New(rand.NewSource(time.Now().UnixNano()))
	charset := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

	b := make([]byte, shortNameLength)
	for i := range b {
		b[i] = charset[randSource.Intn(len(charset))]
	}

	return string(b)
}
