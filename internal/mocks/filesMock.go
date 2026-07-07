package mocks

import (
	"insightly/internal/models"

	"github.com/stretchr/testify/mock"
)

type MockFilesRepo struct {
	mock.Mock
}

func (m *MockFilesRepo) CreateFiles(f models.Files) (models.Files, error) {
	args := m.Called(f)
	return args.Get(0).(models.Files), args.Error(1)
}

func (m *MockFilesRepo) GetFilesByUserId(userId int) ([]models.Files, error) {
	args := m.Called(userId)
	return args.Get(0).([]models.Files), args.Error(1)
}

func (m *MockFilesRepo) GetFileById(id int) (models.Files, error) {
	args := m.Called(id)
	return args.Get(0).(models.Files), args.Error(1)
}

func (m *MockFilesRepo) DeleteFile(id int) error {
	args := m.Called(id)
	return args.Error(0)
}
