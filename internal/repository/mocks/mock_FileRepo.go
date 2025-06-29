// Package mocks - mocks файлового репозитория
package mocks

import (
	"context"

	"github.com/Alheor/shorturl/internal/models"

	"github.com/stretchr/testify/mock"
)

// MockFileRepo - структура файлового репозитория.
type MockFileRepo struct {
	mock.Mock
}

// Add Добавить URL.
func (m *MockFileRepo) Add(ctx context.Context, user *models.User, name string) (string, error) {

	args := m.Called(ctx, user, name)
	return args.String(0), args.Error(0)
}

// AddBatch Добавить несколько URL.
func (m *MockFileRepo) AddBatch(ctx context.Context, user *models.User, list *[]models.BatchEl) error {
	args := m.Called(ctx, user, list)
	return args.Error(0)
}

// GetByShortName Получить URL по короткому имени.
func (m *MockFileRepo) GetByShortName(ctx context.Context, user *models.User, name string) (string, bool, error) {
	args := m.Called(ctx, user, name)
	return args.String(0), args.Bool(1), args.Error(0)
}

// IsReady Готовность репозитория.
func (m *MockFileRepo) IsReady(ctx context.Context) bool {
	args := m.Called(ctx)
	return args.Bool(0)
}

// RemoveByOriginalURL удалить URL.
func (m *MockFileRepo) RemoveByOriginalURL(ctx context.Context, user *models.User, originalURL string) error {
	args := m.Called(ctx, user, originalURL)
	return args.Error(0)
}

// GetAll получить все URL пользователя.
func (m *MockFileRepo) GetAll(ctx context.Context, user *models.User) (<-chan models.HistoryEl, <-chan error) {
	ch := make(chan models.HistoryEl)
	chRrr := make(chan error)
	close(ch)
	close(chRrr)

	return ch, chRrr
}

// RemoveBatch массовое удаление URL.
func (m *MockFileRepo) RemoveBatch(ctx context.Context, user *models.User, list []string) error {
	args := m.Called(ctx, user, list)
	return args.Error(0)
}
