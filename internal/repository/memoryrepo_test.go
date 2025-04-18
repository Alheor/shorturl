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

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMemoryGetUrlNotExists(t *testing.T) {
	err := logger.Init(nil)
	require.NoError(t, err)

	cfg := config.Load()
	cfg.FileStoragePath = ``

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	err = Init(ctx, &cfg, nil)
	require.NoError(t, err)

	url, err := GetRepository().GetByShortName(ctx, user, `any_url`)
	require.NoError(t, err)
	assert.Empty(t, url)
}

func TestMemoryAddURLAndGetURLSuccess(t *testing.T) {
	err := logger.Init(nil)
	require.NoError(t, err)

	cfg := config.Load()
	cfg.FileStoragePath = ``

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	err = Init(ctx, &cfg, nil)
	require.NoError(t, err)

	urlList := map[int]string{1: targetURL + `1`, 2: targetURL + `2`, 3: targetURL + `3`}
	shortsList := make(map[string]string)

	for _, val := range urlList {
		hash, err := GetRepository().Add(ctx, user, val)
		require.NoError(t, err)

		shortsList[hash] = val
	}

	for index, val := range shortsList {
		res, err := GetRepository().GetByShortName(ctx, user, index)
		require.NoError(t, err)
		assert.Equal(t, val, res)
	}
}

func TestMemoryAddURLAndGetAllURLSuccess(t *testing.T) {
	err := logger.Init(nil)
	require.NoError(t, err)

	cfg := config.Load()
	cfg.FileStoragePath = ``

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	err = Init(ctx, &cfg, nil)
	require.NoError(t, err)

	urlList := map[int]string{1: targetURL + `1`, 2: targetURL + `2`, 3: targetURL + `3`}
	shortsList := make(map[string]string)

	for _, val := range urlList {
		hash, err := GetRepository().Add(ctx, user, val)
		require.NoError(t, err)

		shortsList[hash] = val
	}

	list, err := GetRepository().GetAll(ctx, user)
	require.NoError(t, err)

	storageList := *list
	for sourceHash, _ := range shortsList {
		_, exists := storageList[sourceHash]
		assert.True(t, exists)

	}
}

func TestMemoryAddExistsURLSuccess(t *testing.T) {
	err := logger.Init(nil)
	require.NoError(t, err)

	cfg := config.Load()
	cfg.FileStoragePath = ``

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	_ = os.Remove(cfg.FileStoragePath)

	err = Init(ctx, &cfg, nil)
	require.NoError(t, err)

	hash, err := GetRepository().Add(ctx, user, targetURL)
	require.NoError(t, err)

	_, err = GetRepository().Add(ctx, user, targetURL)
	var uniqError *models.UniqueErr
	require.ErrorAs(t, err, &uniqError)
	require.Equal(t, hash, uniqError.ShortKey)
}

func TestMemoryAddBatchSuccess(t *testing.T) {
	err := logger.Init(nil)
	require.NoError(t, err)

	cfg := config.Load()
	cfg.FileStoragePath = ``

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	err = Init(ctx, &cfg, nil)
	require.NoError(t, err)

	var urlList []models.BatchEl

	urlList = append(urlList, models.BatchEl{CorrelationID: `1`, OriginalURL: targetURL + `1`, ShortURL: `hash1`})
	urlList = append(urlList, models.BatchEl{CorrelationID: `2`, OriginalURL: targetURL + `2`, ShortURL: `hash2`})
	urlList = append(urlList, models.BatchEl{CorrelationID: `3`, OriginalURL: targetURL + `3`, ShortURL: `hash3`})

	err = GetRepository().AddBatch(ctx, user, &urlList)
	require.NoError(t, err)

	for _, v := range urlList {
		res, err := GetRepository().GetByShortName(ctx, user, v.ShortURL)
		require.NoError(t, err)
		assert.Equal(t, v.OriginalURL, res)
	}
}

func TestMemoryIsReadySuccess(t *testing.T) {
	cfg := config.Load()
	cfg.FileStoragePath = ``

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	_ = os.Remove(cfg.FileStoragePath)

	mockRepo := new(mocks.MockMemoryRepo)
	mockRepo.On("IsReady", ctx).Return(true)

	err := Init(ctx, &cfg, mockRepo)
	require.NoError(t, err)

	assert.True(t, GetRepository().IsReady(ctx))
}

func TestMemoryIsReadyFileFalse(t *testing.T) {

	cfg := config.Load()
	cfg.FileStoragePath = ``

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	_ = os.Remove(cfg.FileStoragePath)

	mockRepo := new(mocks.MockMemoryRepo)
	mockRepo.On("IsReady", ctx).Return(false)

	err := Init(ctx, &cfg, mockRepo)
	require.NoError(t, err)

	assert.False(t, GetRepository().IsReady(ctx))
}
