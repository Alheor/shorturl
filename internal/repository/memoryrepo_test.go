package repository

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/Alheor/shorturl/internal/config"
	"github.com/Alheor/shorturl/internal/logger"
	"github.com/Alheor/shorturl/internal/models"
	"github.com/Alheor/shorturl/internal/repository/mocks"
	"github.com/Alheor/shorturl/internal/urlhasher"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMemoryGetUrlNotExists(t *testing.T) {
	err := logger.Init(nil)
	require.NoError(t, err)

	cfg := config.Options{FileStoragePath: ``}
	config.Load(&cfg)

	urlhasher.Init(nil)

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	err = Init(ctx, nil)
	require.NoError(t, err)

	url, err := GetRepository().GetByShortName(ctx, `any_url`)
	require.NoError(t, err)
	assert.Empty(t, url)
}

func TestMemoryAddURLAndGetURLSuccess(t *testing.T) {
	err := logger.Init(nil)
	require.NoError(t, err)

	cfg := config.Options{FileStoragePath: ``}
	config.Load(&cfg)

	urlhasher.Init(nil)

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

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
}

func TestMemoryAddExistsURLSuccess(t *testing.T) {
	err := logger.Init(nil)
	require.NoError(t, err)

	cfg := config.Options{FileStoragePath: ``}
	config.Load(&cfg)

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
}

func TestMemoryAddBatchSuccess(t *testing.T) {
	err := logger.Init(nil)
	require.NoError(t, err)

	cfg := config.Options{FileStoragePath: ``}
	config.Load(&cfg)

	urlhasher.Init(nil)

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	err = Init(ctx, nil)
	require.NoError(t, err)

	var urlList []models.BatchEl

	urlList = append(urlList, models.BatchEl{CorrelationID: `1`, OriginalURL: targetURL + `1`, ShortURL: `hash1`})
	urlList = append(urlList, models.BatchEl{CorrelationID: `2`, OriginalURL: targetURL + `2`, ShortURL: `hash2`})
	urlList = append(urlList, models.BatchEl{CorrelationID: `3`, OriginalURL: targetURL + `3`, ShortURL: `hash3`})

	err = GetRepository().AddBatch(ctx, &urlList)
	require.NoError(t, err)

	for _, v := range urlList {
		res, err := GetRepository().GetByShortName(ctx, v.ShortURL)
		require.NoError(t, err)
		assert.Equal(t, v.OriginalURL, res)
	}
}

func TestMemoryIsReadySuccess(t *testing.T) {
	urlhasher.Init(nil)

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	_ = os.Remove(config.GetOptions().FileStoragePath)

	mockRepo := new(mocks.MockMemoryRepo)
	mockRepo.On("IsReady", ctx).Return(true)

	err := Init(ctx, mockRepo)
	require.NoError(t, err)

	assert.True(t, GetRepository().IsReady(ctx))
}

func TestMemoryIsReadyFileFalse(t *testing.T) {
	urlhasher.Init(nil)

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	_ = os.Remove(config.GetOptions().FileStoragePath)

	mockRepo := new(mocks.MockMemoryRepo)
	mockRepo.On("IsReady", ctx).Return(false)

	err := Init(ctx, mockRepo)
	require.NoError(t, err)

	assert.False(t, GetRepository().IsReady(ctx))
}
