package repository

import (
	"os"
	"testing"

	"github.com/Alheor/shorturl/internal/config"
	"github.com/Alheor/shorturl/internal/urlhasher"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const targetURL = `https://practicum.yandex.ru/`

func TestGetUrlNotExists(t *testing.T) {
	config.Load()
	urlhasher.Init()

	_ = os.Remove(config.GetOptions().FileStoragePath)

	err := Init()
	require.NoError(t, err)

	url := GetRepository().GetByShortName(`any_url`)
	assert.Nil(t, url)

	err = os.Remove(config.GetOptions().FileStoragePath)
	require.NoError(t, err)
}

func TestAddURLAndGetURLSuccess(t *testing.T) {
	config.Load()
	urlhasher.Init()

	_ = os.Remove(config.GetOptions().FileStoragePath)

	err := Init()
	require.NoError(t, err)

	urlList := map[int]string{1: targetURL + `1`, 2: targetURL + `2`, 3: targetURL + `3`}
	shortsList := make(map[string]string)

	for _, val := range urlList {
		hash, err := GetRepository().Add(val)
		require.NoError(t, err)

		shortsList[*hash] = val
	}

	for index, val := range shortsList {
		res := GetRepository().GetByShortName(index)
		require.NoError(t, err)
		assert.Equal(t, val, *res)
	}

	err = os.Remove(config.GetOptions().FileStoragePath)
	require.NoError(t, err)
}

func TestAddExistsURLFileSuccess(t *testing.T) {
	config.Load()
	urlhasher.Init()

	_ = os.Remove(config.GetOptions().FileStoragePath)

	err := Init()
	require.NoError(t, err)

	hash, err := GetRepository().Add(targetURL)
	require.NoError(t, err)

	hash1, err := GetRepository().Add(targetURL)
	require.NoError(t, err)
	require.Equal(t, hash, hash1)

	err = os.Remove(config.GetOptions().FileStoragePath)
	require.NoError(t, err)
}

func TestCreatedFileSuccess(t *testing.T) {
	config.Load()
	urlhasher.Init()

	_ = os.Remove(config.GetOptions().FileStoragePath)

	err := Init()
	require.NoError(t, err)

	_, err = GetRepository().Add(targetURL)
	require.NoError(t, err)

	assert.FileExists(t, config.GetOptions().FileStoragePath)

	err = os.Remove(config.GetOptions().FileStoragePath)
	require.NoError(t, err)
}

func TestLoadFromFileSuccess(t *testing.T) {
	config.Load()
	urlhasher.Init()

	_ = os.Remove(config.GetOptions().FileStoragePath)

	err := Init()
	require.NoError(t, err)

	hash, err := GetRepository().Add(targetURL)
	require.NoError(t, err)

	err = Init()
	require.NoError(t, err)

	url := GetRepository().GetByShortName(*hash)
	assert.Equal(t, targetURL, *url)

	err = os.Remove(config.GetOptions().FileStoragePath)
	require.NoError(t, err)
}

func TestIsReadyFileSuccess(t *testing.T) {
	config.Load()
	urlhasher.Init()

	_ = os.Remove(config.GetOptions().FileStoragePath)

	err := Init()
	require.NoError(t, err)

	assert.True(t, GetRepository().IsReady())
}
