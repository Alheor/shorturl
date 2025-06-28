package service

import (
	"context"
	"testing"
	"time"

	"github.com/Alheor/shorturl/internal/config"
	"github.com/Alheor/shorturl/internal/logger"
	"github.com/Alheor/shorturl/internal/models"
	"github.com/Alheor/shorturl/internal/repository"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAddSuccess(t *testing.T) {
	cfg := config.Load()
	cfg.FileStoragePath = ``

	err := logger.Init(nil)
	require.NoError(t, err)

	ctx := context.Background()

	err = repository.Init(ctx, &cfg, nil)
	require.NoError(t, err)

	Init(&cfg)

	user := &models.User{ID: `6a30af51-b6ac-63ba-9e1c-5da06e1b610e`}

	shortURL, err := Add(ctx, user, `https://example.com/?var1=value1&var2=value2`)
	require.NoError(t, err)

	assert.NotEmpty(t, shortURL)
}

func TestAddError(t *testing.T) {
	cfg := config.Load()
	cfg.FileStoragePath = ``

	err := logger.Init(nil)
	require.NoError(t, err)

	ctx := context.Background()

	err = repository.Init(ctx, &cfg, nil)
	require.NoError(t, err)

	Init(&cfg)

	user := &models.User{ID: `6a30af51-b6ac-63ba-9e1c-5da06e1b610e`}

	shortURL, err := Add(ctx, user, `https://example.com/?var1=value1&var2=value2`)
	require.NoError(t, err)

	assert.NotEmpty(t, shortURL)

	_, err = Add(ctx, user, `https://example.com/?var1=value1&var2=value2`)
	require.Error(t, err)
}

func TestGetSuccess(t *testing.T) {
	cfg := config.Load()
	cfg.FileStoragePath = ``

	err := logger.Init(nil)
	require.NoError(t, err)

	ctx := context.Background()

	err = repository.Init(ctx, &cfg, nil)
	require.NoError(t, err)

	Init(&cfg)

	user := &models.User{ID: `6a30af51-b6ac-63ba-9e1c-5da06e1b610e`}
	originalURL := "https://example.com/?var1=value1&var2=value2"

	shortURL, err := Add(ctx, user, originalURL)
	require.NoError(t, err)

	URL, isRemoved := Get(ctx, user, shortURL)
	require.NoError(t, err)

	assert.False(t, isRemoved)
	assert.Equal(t, originalURL, URL)
}

func TestGetError(t *testing.T) {
	cfg := config.Load()
	cfg.FileStoragePath = ``

	err := logger.Init(nil)
	require.NoError(t, err)

	ctx := context.Background()

	err = repository.Init(ctx, &cfg, nil)
	require.NoError(t, err)

	Init(&cfg)

	user := &models.User{ID: `6a30af51-b6ac-63ba-9e1c-5da06e1b610e`}

	URL, isRemoved := Get(ctx, user, `short_name`)

	assert.False(t, isRemoved)
	assert.Empty(t, URL)
}

func TestAddBatchSuccess(t *testing.T) {
	cfg := config.Load()
	cfg.FileStoragePath = ``

	err := logger.Init(nil)
	require.NoError(t, err)

	ctx := context.Background()

	err = repository.Init(ctx, &cfg, nil)
	require.NoError(t, err)

	Init(&cfg)

	user := &models.User{ID: `6a30af51-b6ac-63ba-9e1c-5da06e1b610e`}

	var list []models.APIBatchRequestEl
	var res []models.APIBatchResponseEl

	list = append(list, models.APIBatchRequestEl{CorrelationID: `1`, OriginalURL: `https://example.com/?var1=value1&var2=value2`})
	list = append(list, models.APIBatchRequestEl{CorrelationID: `2`, OriginalURL: `https://example.com/?var2=value2&var3=value3`})
	list = append(list, models.APIBatchRequestEl{CorrelationID: `3`, OriginalURL: `https://example.com/?var3=value3&var4=value4`})

	res, err = AddBatch(ctx, user, list)
	require.NoError(t, err)

	assert.Len(t, res, 3)
}

func TestAddBatchError(t *testing.T) {
	cfg := config.Load()
	cfg.FileStoragePath = ``

	err := logger.Init(nil)
	require.NoError(t, err)

	ctx := context.Background()

	err = repository.Init(ctx, &cfg, nil)
	require.NoError(t, err)

	Init(&cfg)

	user := &models.User{ID: `6a30af51-b6ac-63ba-9e1c-5da06e1b610e`}

	var list []models.APIBatchRequestEl

	list = append(list, models.APIBatchRequestEl{CorrelationID: `1`, OriginalURL: `https://example.com/?var1=value1&var2=value2`})
	list = append(list, models.APIBatchRequestEl{CorrelationID: `2`, OriginalURL: `https://example.com/?var2=value2&var3=value3`})
	list = append(list, models.APIBatchRequestEl{CorrelationID: `3`, OriginalURL: `https://example.com/?var3=value3&var4=value4`})

	ctx, cancel := context.WithTimeout(ctx, 10*time.Millisecond)
	defer cancel()

	time.Sleep(11 * time.Millisecond)

	_, err = AddBatch(ctx, user, list)
	require.Error(t, err)
}

func TestGetAllSuccess(t *testing.T) {
	cfg := config.Load()
	cfg.FileStoragePath = ``

	err := logger.Init(nil)
	require.NoError(t, err)

	ctx := context.Background()

	err = repository.Init(ctx, &cfg, nil)
	require.NoError(t, err)

	Init(&cfg)

	user := &models.User{ID: `6a30af51-b6ac-63ba-9e1c-5da06e1b610e`}

	var list []models.APIBatchRequestEl
	var res []models.APIBatchResponseEl

	list = append(list, models.APIBatchRequestEl{CorrelationID: `1`, OriginalURL: `https://example.com/?var1=value1&var2=value2`})
	list = append(list, models.APIBatchRequestEl{CorrelationID: `2`, OriginalURL: `https://example.com/?var2=value2&var3=value3`})
	list = append(list, models.APIBatchRequestEl{CorrelationID: `3`, OriginalURL: `https://example.com/?var3=value3&var4=value4`})

	res, err = AddBatch(ctx, user, list)
	require.NoError(t, err)
	assert.Len(t, res, 3)

	var allList []models.HistoryEl
	var errList []error

	chList, chErr := GetAll(ctx, user)

	for el := range chList {
		allList = append(allList, el)
	}

	for sss := range chErr {
		errList = append(errList, sss)
	}

	assert.Len(t, allList, 3)
	assert.Len(t, errList, 0)
}

func TestGetAllError(t *testing.T) {
	cfg := config.Load()
	cfg.FileStoragePath = ``

	err := logger.Init(nil)
	require.NoError(t, err)

	ctx := context.Background()

	err = repository.Init(ctx, &cfg, nil)
	require.NoError(t, err)

	Init(&cfg)

	user := &models.User{ID: `6a30af51-b6ac-63ba-9e1c-5da06e1b610e`}

	var list []models.APIBatchRequestEl
	var res []models.APIBatchResponseEl

	list = append(list, models.APIBatchRequestEl{CorrelationID: `1`, OriginalURL: `https://example.com/?var1=value1&var2=value2`})
	list = append(list, models.APIBatchRequestEl{CorrelationID: `2`, OriginalURL: `https://example.com/?var2=value2&var3=value3`})
	list = append(list, models.APIBatchRequestEl{CorrelationID: `3`, OriginalURL: `https://example.com/?var3=value3&var4=value4`})

	res, err = AddBatch(ctx, user, list)
	require.NoError(t, err)
	assert.Len(t, res, 3)

	var allList []models.HistoryEl
	var errList []error

	ctx, cancel := context.WithTimeout(ctx, 10*time.Millisecond)
	defer cancel()

	time.Sleep(11 * time.Millisecond)

	chList, chErr := GetAll(ctx, user)

	for el := range chList {
		allList = append(allList, el)
	}

	for sss := range chErr {
		errList = append(errList, sss)
	}

	assert.Len(t, allList, 0)
	assert.Len(t, errList, 1)
}

func TestRemoveBatchError(t *testing.T) {
	cfg := config.Load()
	cfg.FileStoragePath = ``

	err := logger.Init(nil)
	require.NoError(t, err)

	ctx := context.Background()

	err = repository.Init(ctx, &cfg, nil)
	require.NoError(t, err)

	Init(&cfg)

	user := &models.User{ID: `6a30af51-b6ac-63ba-9e1c-5da06e1b610e`}

	var list []string
	list = append(list, `short_name1`)
	list = append(list, `short_name2`)
	list = append(list, `short_name3`)

	err = RemoveBatch(ctx, user, list)
	require.Error(t, err)
}

func TestIsDBReadySuccess(t *testing.T) {
	cfg := config.Load()
	cfg.FileStoragePath = ``

	err := logger.Init(nil)
	require.NoError(t, err)

	ctx := context.Background()

	err = repository.Init(ctx, &cfg, nil)
	require.NoError(t, err)

	Init(&cfg)

	isReady := IsDBReady(ctx)
	require.True(t, isReady)
}
