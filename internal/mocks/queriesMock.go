package mocks

import (
	"insightly/internal/models"

	"github.com/stretchr/testify/mock"
)

type MockQueriesRepo struct {
	mock.Mock
}

func (m *MockQueriesRepo) CreateQueries(q models.Queries) (models.Queries, error) {
	args := m.Called(q)
	return args.Get(0).(models.Queries), args.Error(1)
}

func (m *MockQueriesRepo) GetQueriesByUserId(id int) ([]models.Queries, error) {
	args := m.Called(id)
	return args.Get(0).([]models.Queries), args.Error(1)
}
