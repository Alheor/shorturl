package repository

import (
	"context"
	"reflect"
	"testing"
	"time"

	"github.com/Alheor/shorturl/internal/config"
	"github.com/Alheor/shorturl/internal/logger"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRepoRunInMemoryMode(t *testing.T) {
	cfg := config.Options{FileStoragePath: ``}
	config.Load(&cfg)

	err := logger.Init(nil)
	require.NoError(t, err)

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	err = Init(ctx, nil)
	require.NoError(t, err)
	assert.Equal(t, `MemoryRepo`, reflect.TypeOf(GetRepository()).Elem().Name())
}

func TestRepoRunInFileMode(t *testing.T) {
	cfg := config.Options{FileStoragePath: `/tmp/short-url.json`}
	config.Load(&cfg)

	err := logger.Init(nil)
	require.NoError(t, err)

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	err = Init(ctx, nil)
	require.NoError(t, err)
	assert.Equal(t, `FileRepo`, reflect.TypeOf(GetRepository()).Elem().Name())
}

func TestRepoRunInDBMode(t *testing.T) {
	cfg := config.Options{DatabaseDsn: `user=app password=pass host=localhost port=5432 dbname=app pool_max_conns=10`}
	config.Load(&cfg)

	err := logger.Init(nil)
	require.NoError(t, err)

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	err = Init(ctx, nil)
	require.NoError(t, err)
	assert.Equal(t, `PostgresRepo`, reflect.TypeOf(GetRepository()).Elem().Name())
}
