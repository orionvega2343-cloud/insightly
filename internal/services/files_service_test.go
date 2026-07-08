package services

import (
	"errors"
	"fmt"
	"insightly/internal/errs"
	"insightly/internal/mocks"
	"insightly/internal/models"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func Test_CreateFilesSuccess(t *testing.T) {
	var (
		userId   = 1
		filename = "test.csv"
		data     []byte
	)
	filesRepo := new(mocks.MockFilesRepo)
	svc := NewFilesService(filesRepo)
	f := models.Files{UserId: userId, Name: filename}

	filesRepo.On("CreateFiles", mock.Anything).Return(f, nil)

	file, err := svc.CreateFiles(userId, filename, data)

	//Очищаем папку после запуска теста
	defer os.RemoveAll(fmt.Sprintf("uploads/%d", userId))

	assert.Equal(t, f.Name, file.Name)
	assert.NoError(t, err)
	filesRepo.AssertExpectations(t)
}

func Test_CreateFilesFail(t *testing.T) {
	var (
		userId   = 1
		filename = "test.csv"
		data     []byte
	)
	filesRepo := new(mocks.MockFilesRepo)
	svc := NewFilesService(filesRepo)
	f := models.Files{UserId: userId, Name: filename}

	filesRepo.On("CreateFiles", mock.Anything).Return(f, errors.New("some error"))

	_, err := svc.CreateFiles(userId, filename, data)

	//Очищаем папку после запуска теста
	defer os.RemoveAll(fmt.Sprintf("uploads/%d", userId))

	assert.Equal(t, errs.FailedCreateFile, err)
	assert.Error(t, err)
	filesRepo.AssertExpectations(t)
}

func Test_GetFilesByUserIdSuccess(t *testing.T) {
	filesRepo := new(mocks.MockFilesRepo)
	svc := NewFilesService(filesRepo)
	f := models.Files{UserId: 1}

	filesRepo.On("GetFilesByUserId", mock.Anything).Return([]models.Files{f}, nil)

	files, err := svc.GetFilesByUserId(1)

	assert.NoError(t, err)
	assert.Equal(t, f.UserId, files[0].UserId)
	filesRepo.AssertExpectations(t)

}

func Test_GetFilesByUserIdFail(t *testing.T) {
	mockErr := errors.New("some error")
	filesRepo := new(mocks.MockFilesRepo)
	svc := NewFilesService(filesRepo)
	f := models.Files{UserId: 1}

	filesRepo.On("GetFilesByUserId", mock.Anything).Return([]models.Files{f}, mockErr)

	_, err := svc.GetFilesByUserId(1)

	assert.Error(t, err)
	assert.Equal(t, mockErr, err)
	filesRepo.AssertExpectations(t)

}

func Test_GetFileByIdSuccess(t *testing.T) {
	filesRepo := new(mocks.MockFilesRepo)
	svc := NewFilesService(filesRepo)
	f := models.Files{UserId: 1}

	filesRepo.On("GetFileById", mock.Anything).Return(f, nil)

	file, err := svc.GetFileById(1)

	assert.Equal(t, f.UserId, file.UserId)
	assert.NoError(t, err)
	filesRepo.AssertExpectations(t)
}

func Test_GetFileByIdFail(t *testing.T) {
	mockErr := errors.New("some error")
	filesRepo := new(mocks.MockFilesRepo)
	svc := NewFilesService(filesRepo)
	f := models.Files{UserId: 1}

	filesRepo.On("GetFileById", mock.Anything).Return(f, mockErr)

	_, err := svc.GetFileById(1)

	assert.Error(t, err)
	assert.Equal(t, mockErr, err)
	filesRepo.AssertExpectations(t)
}

func Test_DeleteFileSuccess(t *testing.T) {
	filesRepo := new(mocks.MockFilesRepo)
	svc := NewFilesService(filesRepo)

	filesRepo.On("DeleteFile", mock.Anything).Return(nil)

	err := svc.DeleteFile(1)

	assert.Equal(t, nil, err)
	assert.NoError(t, err)
	filesRepo.AssertExpectations(t)

}

func Test_DeleteFileFail(t *testing.T) {
	mockErr := errors.New("some error")
	filesRepo := new(mocks.MockFilesRepo)
	svc := NewFilesService(filesRepo)

	filesRepo.On("DeleteFile", mock.Anything).Return(mockErr)

	err := svc.DeleteFile(1)

	assert.Equal(t, mockErr, err)
	assert.Error(t, err)
	filesRepo.AssertExpectations(t)

}
