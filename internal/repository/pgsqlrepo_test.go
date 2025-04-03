package repository

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/Alheor/shorturl/internal/config"
	"github.com/Alheor/shorturl/internal/repository/mocks"
	"github.com/Alheor/shorturl/internal/urlhasher"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestIsReadyDBSuccess(t *testing.T) {
	config.Load()
	urlhasher.Init(nil)

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	_ = os.Remove(config.GetOptions().FileStoragePath)

	mockRepo := new(mocks.MockPostgres)
	mockRepo.On("IsReady", ctx).Return(true)

	err := Init(ctx, mockRepo)
	require.NoError(t, err)

	assert.True(t, GetRepository().IsReady(ctx))
}

func TestIsReadyDBFail(t *testing.T) {
	config.Load()
	urlhasher.Init(nil)

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	_ = os.Remove(config.GetOptions().FileStoragePath)

	mockRepo := new(mocks.MockPostgres)
	mockRepo.On("IsReady", ctx).Return(false)

	err := Init(ctx, mockRepo)
	require.NoError(t, err)

	assert.False(t, GetRepository().IsReady(ctx))
}
