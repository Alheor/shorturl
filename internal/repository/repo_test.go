package repository

import (
	"testing"

	"github.com/Alheor/shorturl/internal/urlhasher"

	"github.com/stretchr/testify/assert"
)

func TestAddUrlWithNewEl(t *testing.T) {
	urlhasher.Init()
	Init()

	hash := Add(`new_url`)
	assert.NotEmpty(t, hash)
}

func TestAddUrlWithExistedEl(t *testing.T) {
	urlhasher.Init()
	Init()

	hash := Add(`new_url`)
	assert.NotEmpty(t, hash)

	hashExists := Add(`new_url`)
	assert.NotEmpty(t, hashExists)
	assert.Equal(t, hash, hashExists)
}

func TestGetUrl(t *testing.T) {
	urlhasher.Init()
	Init()

	hash := Add(`new_url`)
	assert.NotEmpty(t, hash)

	url := GetByShortName(hash)
	assert.NotEmpty(t, url)
	assert.Equal(t, `new_url`, *url)
}
