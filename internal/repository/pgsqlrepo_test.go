package repository

import (
	"context"
	"testing"
	"time"

	"github.com/Alheor/shorturl/internal/repository/mocks"
	"github.com/Alheor/shorturl/internal/urlhasher"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Для ручного запуска с локальной БД
//
//func TestDBGetUrlNotExists(t *testing.T) {
//	err := logger.Init(nil)
//	require.NoError(t, err)
//
//	cfg := config.Options{DatabaseDsn: `user=app password=pass host=localhost port=5432 dbname=app pool_max_conns=10`}
//	config.Load(&cfg)
//
//	urlhasher.Init(nil)
//
//	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
//	defer cancel()
//
//	err = Init(ctx, nil)
//	require.NoError(t, err)
//
//	url, err := GetRepository().GetByShortName(ctx, `any_url`)
//	require.NoError(t, err)
//	assert.Empty(t, url)
//}

// Для ручного запуска с локальной БД
//
//func TestDBAddURLAndGetURLSuccess(t *testing.T) {
//	err := logger.Init(nil)
//	require.NoError(t, err)
//
//	cfg := config.Options{DatabaseDsn: `user=app password=pass host=localhost port=5432 dbname=app pool_max_conns=10`}
//	config.Load(&cfg)
//
//	urlhasher.Init(nil)
//
//	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
//	defer cancel()
//
//	err = Init(ctx, nil)
//	require.NoError(t, err)
//
//	urlList := map[int]string{1: targetURL + `1`, 2: targetURL + `2`, 3: targetURL + `3`}
//	shortsList := make(map[string]string)
//
//	for _, val := range urlList {
//		hash, err := GetRepository().Add(ctx, val)
//		require.NoError(t, err)
//
//		shortsList[hash] = val
//	}
//
//	for index, val := range shortsList {
//		res, err := GetRepository().GetByShortName(ctx, index)
//		require.NoError(t, err)
//		assert.Equal(t, val, res)
//	}
//}

func TestDBIsReadySuccess(t *testing.T) {
	urlhasher.Init(nil)

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	mockRepo := new(mocks.MockPostgres)
	mockRepo.On("IsReady", ctx).Return(true)

	err := Init(ctx, mockRepo)
	require.NoError(t, err)

	assert.True(t, GetRepository().IsReady(ctx))
}

func TestDBIsReadyFail(t *testing.T) {
	urlhasher.Init(nil)

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	mockRepo := new(mocks.MockPostgres)
	mockRepo.On("IsReady", ctx).Return(false)

	err := Init(ctx, mockRepo)
	require.NoError(t, err)

	assert.False(t, GetRepository().IsReady(ctx))
}
