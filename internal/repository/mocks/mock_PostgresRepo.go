package mocks

import (
	"context"

	"github.com/Alheor/shorturl/internal/models"

	"github.com/stretchr/testify/mock"
)

type MockPostgres struct {
	mock.Mock
}

func (m *MockPostgres) Add(ctx context.Context, user *models.User, name string) (string, error) {

	args := m.Called(ctx, user, name)
	return args.String(0), args.Error(0)
}

func (m *MockPostgres) AddBatch(ctx context.Context, user *models.User, list *[]models.BatchEl) error {
	args := m.Called(ctx, list)
	return args.Error(0)
}

func (m *MockPostgres) GetByShortName(ctx context.Context, user *models.User, name string) (string, error) {

	args := m.Called(ctx, user, name)
	return args.String(0), args.Error(0)
}

func (m *MockPostgres) IsReady(ctx context.Context) bool {
	args := m.Called(ctx)
	return args.Bool(0)
}

// RemoveByOriginalURL удалить url
func (m *MockPostgres) RemoveByOriginalURL(ctx context.Context, user *models.User, originalURL string) error {
	args := m.Called(ctx, user, originalURL)
	return args.Error(0)
}

func (m *MockPostgres) GetAll(ctx context.Context, user *models.User) (*map[string]string, error) {
	args := m.Called(ctx, user)
	return &map[string]string{}, args.Error(0)
}
