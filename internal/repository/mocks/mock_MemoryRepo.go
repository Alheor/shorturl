package mocks

import (
	"context"

	"github.com/Alheor/shorturl/internal/models"

	"github.com/stretchr/testify/mock"
)

type MockMemoryRepo struct {
	mock.Mock
}

func (m *MockMemoryRepo) Add(ctx context.Context, name string) (string, error) {

	args := m.Called(ctx, name)
	return args.String(0), args.Error(0)
}

func (m *MockMemoryRepo) AddBatch(ctx context.Context, list *[]models.BatchEl) error {
	args := m.Called(ctx, list)
	return args.Error(0)
}

func (m *MockMemoryRepo) GetByShortName(ctx context.Context, name string) (string, error) {

	args := m.Called(ctx, name)
	return args.String(0), args.Error(1)
}

func (m *MockMemoryRepo) IsReady(ctx context.Context) bool {
	args := m.Called(ctx)
	return args.Bool(0)
}
