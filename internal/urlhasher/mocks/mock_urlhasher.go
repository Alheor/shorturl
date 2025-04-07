package mocks

import (
	"github.com/stretchr/testify/mock"
)

type MockShortName struct {
	mock.Mock
}

func (m *MockShortName) Generate() string {
	args := m.Called()
	return args.String(0)
}
