package mocks

import (
	"insightly/internal/models"

	"github.com/stretchr/testify/mock"
)

type MockTokensRepo struct {
	mock.Mock
}

func (m *MockTokensRepo) CreateToken(t models.RefreshToken) (*models.RefreshToken, error) {
	args := m.Called(t)
	return args.Get(0).(*models.RefreshToken), args.Error(1)
}

func (m *MockTokensRepo) GetTokenByValue(token string) (*models.RefreshToken, error) {
	args := m.Called(token)
	return args.Get(0).(*models.RefreshToken), args.Error(1)
}

func (m *MockTokensRepo) DeleteToken(id int) error {
	args := m.Called(id)
	return args.Error(0)
}
