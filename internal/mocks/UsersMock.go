package mocks

import (
	"insightly/internal/models"

	"github.com/stretchr/testify/mock"
)

type MockUserRepo struct {
	mock.Mock
}

func (m *MockUserRepo) CreateUser(u models.User) (models.User, error) {
	args := m.Called(u)
	return args.Get(0).(models.User), args.Error(1)
}

func (m *MockUserRepo) GetUserByEmail(email string) (models.User, error) {
	args := m.Called(email)
	return args.Get(0).(models.User), args.Error(1)
}
