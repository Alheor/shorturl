package mocks

import (
	"context"

	"github.com/Alheor/shorturl/internal/models"

	"github.com/stretchr/testify/mock"
)

type MockPostgres struct {
	mock.Mock
}

func (m *MockPostgres) Add(ctx context.Context, name string) (string, error) {

	args := m.Called(ctx, name)
	return args.String(0), args.Error(0)
}

func (m *MockPostgres) AddBatch(ctx context.Context, list *[]models.BatchEl) error {
	args := m.Called(ctx, list)
	return args.Error(0)
}

func (m *MockPostgres) GetByShortName(ctx context.Context, name string) (string, error) {

	args := m.Called(ctx, name)
	return args.String(0), args.Error(0)
}

func (m *MockPostgres) IsReady(ctx context.Context) bool {
	args := m.Called(ctx)
	return args.Bool(0)
}
