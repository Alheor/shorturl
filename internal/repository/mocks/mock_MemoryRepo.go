package mocks

import (
	"context"

	"github.com/Alheor/shorturl/internal/models"

	"github.com/stretchr/testify/mock"
)

type MockMemoryRepo struct {
	mock.Mock
}

func (m *MockMemoryRepo) Add(ctx context.Context, user *models.User, name string) (string, error) {

	args := m.Called(ctx, user, name)
	return args.String(0), args.Error(0)
}

func (m *MockMemoryRepo) AddBatch(ctx context.Context, user *models.User, list *[]models.BatchEl) error {
	args := m.Called(ctx, user, list)
	return args.Error(0)
}

func (m *MockMemoryRepo) GetByShortName(ctx context.Context, user *models.User, name string) (string, error) {

	args := m.Called(ctx, user, name)
	return args.String(0), args.Error(1)
}

func (m *MockMemoryRepo) IsReady(ctx context.Context) bool {
	args := m.Called(ctx)
	return args.Bool(0)
}

// RemoveByOriginalURL удалить url
func (m *MockMemoryRepo) RemoveByOriginalURL(ctx context.Context, user *models.User, originalURL string) error {
	args := m.Called(ctx, user, originalURL)
	return args.Error(0)
}

func (m *MockMemoryRepo) GetAll(ctx context.Context, user *models.User) (*map[string]string, error) {
	args := m.Called(ctx, user)
	return &map[string]string{}, args.Error(0)
}
