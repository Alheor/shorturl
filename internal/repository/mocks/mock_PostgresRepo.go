// Package mocks - mocks репозитория PostgreSQl
package mocks

import (
	"context"

	"github.com/Alheor/shorturl/internal/models"

	"github.com/stretchr/testify/mock"
)

// MockPostgres - структура файлового репозитория.
type MockPostgres struct {
	mock.Mock
}

// Add Добавить URL.
func (m *MockPostgres) Add(ctx context.Context, user *models.User, name string) (string, error) {

	args := m.Called(ctx, user, name)
	return args.String(0), args.Error(0)
}

// AddBatch Добавить несколько URL.
func (m *MockPostgres) AddBatch(ctx context.Context, user *models.User, list *[]models.BatchEl) error {
	args := m.Called(ctx, list)
	return args.Error(0)
}

// GetByShortName Получить URL по короткому имени.
func (m *MockPostgres) GetByShortName(ctx context.Context, user *models.User, name string) (string, bool, error) {

	args := m.Called(ctx, user, name)
	return args.String(0), args.Bool(1), args.Error(0)
}

// IsReady Готовность репозитория.
func (m *MockPostgres) IsReady(ctx context.Context) bool {
	args := m.Called(ctx)
	return args.Bool(0)
}

// RemoveByOriginalURL удалить URL.
func (m *MockPostgres) RemoveByOriginalURL(ctx context.Context, user *models.User, originalURL string) error {
	args := m.Called(ctx, user, originalURL)
	return args.Error(0)
}

// GetAll получить все URL пользователя.
func (m *MockPostgres) GetAll(ctx context.Context, user *models.User) (<-chan models.HistoryEl, <-chan error) {
	ch := make(chan models.HistoryEl)
	chRrr := make(chan error)
	close(ch)
	close(chRrr)

	return ch, chRrr
}

// RemoveBatch массовое удаление URL.
func (m *MockPostgres) RemoveBatch(ctx context.Context, user *models.User, list []string) error {
	args := m.Called(ctx, user, list)
	return args.Error(0)
}
