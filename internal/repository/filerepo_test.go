package repository

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/Alheor/shorturl/internal/config"
	"github.com/Alheor/shorturl/internal/logger"
	"github.com/Alheor/shorturl/internal/repository/mocks"
	"github.com/Alheor/shorturl/internal/urlhasher"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const targetURL = `https://practicum.yandex.ru/`

func TestFileGetUrlNotExists(t *testing.T) {
	err := logger.Init(nil)
	require.NoError(t, err)

	config.Load(nil)
	urlhasher.Init(nil)

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	_ = os.Remove(config.GetOptions().FileStoragePath)

	err = Init(ctx, nil)
	require.NoError(t, err)

	url, err := GetRepository().GetByShortName(ctx, `any_url`)
	require.NoError(t, err)
	assert.Empty(t, url)

	err = os.Remove(config.GetOptions().FileStoragePath)
	require.NoError(t, err)
}

func TestFileAddURLAndGetURLSuccess(t *testing.T) {
	err := logger.Init(nil)
	require.NoError(t, err)

	urlhasher.Init(nil)

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	_ = os.Remove(config.GetOptions().FileStoragePath)

	err = Init(ctx, nil)
	require.NoError(t, err)

	urlList := map[int]string{1: targetURL + `1`, 2: targetURL + `2`, 3: targetURL + `3`}
	shortsList := make(map[string]string)

	for _, val := range urlList {
		hash, err := GetRepository().Add(ctx, val)
		require.NoError(t, err)

		shortsList[hash] = val
	}

	for index, val := range shortsList {
		res, err := GetRepository().GetByShortName(ctx, index)
		require.NoError(t, err)
		assert.Equal(t, val, res)
	}

	err = os.Remove(config.GetOptions().FileStoragePath)
	require.NoError(t, err)
}

func TestFileAddExistsURLFileSuccess(t *testing.T) {
	err := logger.Init(nil)
	require.NoError(t, err)

	urlhasher.Init(nil)

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	_ = os.Remove(config.GetOptions().FileStoragePath)

	err = Init(ctx, nil)
	require.NoError(t, err)

	hash, err := GetRepository().Add(ctx, targetURL)
	require.NoError(t, err)

	hash1, err := GetRepository().Add(ctx, targetURL)
	require.NoError(t, err)
	require.Equal(t, hash, hash1)

	err = os.Remove(config.GetOptions().FileStoragePath)
	require.NoError(t, err)
}

func TestFileCreatedFileSuccess(t *testing.T) {
	err := logger.Init(nil)
	require.NoError(t, err)

	urlhasher.Init(nil)

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	_ = os.Remove(config.GetOptions().FileStoragePath)

	err = Init(ctx, nil)
	require.NoError(t, err)

	_, err = GetRepository().Add(ctx, targetURL)
	require.NoError(t, err)

	assert.FileExists(t, config.GetOptions().FileStoragePath)

	err = os.Remove(config.GetOptions().FileStoragePath)
	require.NoError(t, err)
}

func TestFileLoadFromFileSuccess(t *testing.T) {
	err := logger.Init(nil)
	require.NoError(t, err)

	urlhasher.Init(nil)

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	_ = os.Remove(config.GetOptions().FileStoragePath)

	err = Init(ctx, nil)
	require.NoError(t, err)

	hash, err := GetRepository().Add(ctx, targetURL)
	require.NoError(t, err)

	err = Init(ctx, nil)
	require.NoError(t, err)

	url, err := GetRepository().GetByShortName(ctx, hash)
	require.NoError(t, err)
	assert.Equal(t, targetURL, url)

	err = os.Remove(config.GetOptions().FileStoragePath)
	require.NoError(t, err)
}

func TestFileIsReadyFileSuccess(t *testing.T) {
	urlhasher.Init(nil)

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	_ = os.Remove(config.GetOptions().FileStoragePath)

	mockRepo := new(mocks.MockFileRepo)
	mockRepo.On("IsReady", ctx).Return(true)

	err := Init(ctx, mockRepo)
	require.NoError(t, err)

	assert.True(t, GetRepository().IsReady(ctx))
}

func TestFileIsReadyFileFalse(t *testing.T) {
	urlhasher.Init(nil)

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	_ = os.Remove(config.GetOptions().FileStoragePath)

	mockRepo := new(mocks.MockFileRepo)
	mockRepo.On("IsReady", ctx).Return(false)

	err := Init(ctx, mockRepo)
	require.NoError(t, err)

	assert.False(t, GetRepository().IsReady(ctx))
}
