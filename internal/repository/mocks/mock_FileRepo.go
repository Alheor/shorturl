package mocks

import (
	"context"

	"github.com/Alheor/shorturl/internal/models"

	"github.com/stretchr/testify/mock"
)

type MockFileRepo struct {
	mock.Mock
}

func (m *MockFileRepo) Add(ctx context.Context, name string) (string, error) {

	args := m.Called(ctx, name)
	return args.String(0), args.Error(0)
}

func (m *MockFileRepo) AddBatch(ctx context.Context, list *[]models.BatchEl) error {
	args := m.Called(ctx, list)
	return args.Error(0)
}

func (m *MockFileRepo) GetByShortName(ctx context.Context, name string) (string, error) {

	args := m.Called(ctx, name)
	return args.String(0), args.Error(0)
}

func (m *MockFileRepo) IsReady(ctx context.Context) bool {
	args := m.Called(ctx)
	return args.Bool(0)
}

// RemoveByOriginalURL удалить url
func (m *MockFileRepo) RemoveByOriginalURL(ctx context.Context, originalURL string) error {
	args := m.Called(ctx, originalURL)
	return args.Error(0)
}
