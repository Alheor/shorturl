package repository

import (
	"context"
	"testing"
	"time"

	"github.com/Alheor/shorturl/internal/config"
	"github.com/Alheor/shorturl/internal/logger"
	"github.com/Alheor/shorturl/internal/models"
	"github.com/Alheor/shorturl/internal/repository/mocks"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDBGetUrlNotExists(t *testing.T) {

	t.Skip(`Run with database only`) // Для ручного запуска с локальной БД

	err := logger.Init(nil)
	require.NoError(t, err)

	cfg := config.Load()
	cfg.DatabaseDsn = `user=app password=pass host=localhost port=5432 dbname=app pool_max_conns=10`

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	err = Init(ctx, &cfg, nil)
	require.NoError(t, err)

	_, err = Connection.Exec(ctx, `TRUNCATE short_url`)
	require.NoError(t, err)

	url, err := GetRepository().GetByShortName(ctx, user, `any_url`)
	require.NoError(t, err)
	assert.Empty(t, url)
}

func TestDBAddURLAndGetURLSuccess(t *testing.T) {

	t.Skip(`Run with database only`) // Для ручного запуска с локальной БД

	err := logger.Init(nil)
	require.NoError(t, err)

	cfg := config.Load()
	cfg.DatabaseDsn = `user=app password=pass host=localhost port=5432 dbname=app pool_max_conns=10`

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	err = Init(ctx, &cfg, nil)
	require.NoError(t, err)

	_, err = Connection.Exec(ctx, `TRUNCATE short_url`)
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

func TestDBAddBatchSuccess(t *testing.T) {

	t.Skip(`Run with database only`) // Для ручного запуска с локальной БД

	err := logger.Init(nil)
	require.NoError(t, err)

	cfg := config.Load()
	cfg.DatabaseDsn = `user=app password=pass host=localhost port=5432 dbname=app pool_max_conns=10`

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	err = Init(ctx, &cfg, nil)
	require.NoError(t, err)

	_, err = Connection.Exec(ctx, `TRUNCATE short_url`)
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

func TestDBIsReadySuccess(t *testing.T) {

	t.Skip(`Run with database only`) // Для ручного запуска с локальной БД

	err := logger.Init(nil)
	require.NoError(t, err)

	cfg := config.Load()
	cfg.DatabaseDsn = `user=app password=pass host=localhost port=5432 dbname=app pool_max_conns=10`

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	err = Init(ctx, &cfg, nil)
	require.NoError(t, err)

	_, err = Connection.Exec(ctx, `TRUNCATE short_url`)
	require.NoError(t, err)

	assert.True(t, GetRepository().IsReady(ctx))
}

func TestDBIsReadyFail(t *testing.T) {
	cfg := config.Load()

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	mockRepo := new(mocks.MockPostgres)
	mockRepo.On("IsReady", ctx).Return(false)

	err := Init(ctx, &cfg, mockRepo)
	require.NoError(t, err)

	assert.False(t, GetRepository().IsReady(ctx))
}
