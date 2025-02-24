package repository

import (
	"testing"

	"github.com/Alheor/shorturl/internal/urlhasher"

	"github.com/stretchr/testify/assert"
)

func TestAddUrlWithNewEl(t *testing.T) {
	urlhasher.Init()

	hash := GetRepository().Add(`new_url`)
	assert.NotEmpty(t, hash)
}

func TestAddUrlWithExistedEl(t *testing.T) {
	urlhasher.Init()

	hash := GetRepository().Add(`new_url`)
	assert.NotEmpty(t, hash)

	hashExists := GetRepository().Add(`new_url`)
	assert.NotEmpty(t, hashExists)
	assert.Equal(t, hash, hashExists)
}

func TestGetUrlExists(t *testing.T) {
	urlhasher.Init()

	hash := GetRepository().Add(`new_url`)
	assert.NotEmpty(t, hash)

	url := GetRepository().GetByShortName(hash)
	assert.NotEmpty(t, url)
	assert.Equal(t, `new_url`, *url)
}

func TestGetUrlNotExists(t *testing.T) {
	urlhasher.Init()

	url := GetRepository().GetByShortName(`any_url`)
	assert.Nil(t, url)
}
