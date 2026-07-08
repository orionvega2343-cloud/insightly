package services

import (
	"errors"
	"insightly/internal/errs"
	"insightly/internal/mocks"
	"insightly/internal/models"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func Test_CreateQueriesSuccess(t *testing.T) {
	var (
		fileId = 1
		userId = 1
		prompt = "Привет,я функция теста системы запросов!"
	)
	filesRepo := new(mocks.MockFilesRepo)
	queriesRepo := new(mocks.MockQueriesRepo)
	ai := new(mocks.AiMock)
	svc := NewQueriesService(filesRepo, queriesRepo, ai)
	q := models.Queries{FileId: fileId, UserId: userId, Question: prompt}

	queriesRepo.On("CreateQueries", mock.Anything).Return(q, nil)
	ai.On("Analyze", mock.Anything, mock.Anything).Return("", nil)
	filesRepo.On("GetFileById", mock.Anything).Return(models.Files{Path: "" + "testdata/test.csv"}, nil)

	query, err := svc.CreateQueries(fileId, userId, prompt)

	assert.Equal(t, fileId, query.FileId)
	assert.NoError(t, err)
	filesRepo.AssertExpectations(t)
	queriesRepo.AssertExpectations(t)
	ai.AssertExpectations(t)
}

func Test_CreateQueriesParserFail(t *testing.T) {
	var (
		fileId  = 1
		userId  = 1
		prompt  = "Привет,я функция теста системы запросов!"
		mockErr = errors.New("mock error")
	)
	filesRepo := new(mocks.MockFilesRepo)
	queriesRepo := new(mocks.MockQueriesRepo)
	ai := new(mocks.AiMock)
	svc := NewQueriesService(filesRepo, queriesRepo, ai)

	filesRepo.On("GetFileById", mock.Anything).Return(models.Files{Path: "" + "testdata/test.csv"}, mockErr)

	_, err := svc.CreateQueries(fileId, userId, prompt)

	assert.Equal(t, errs.ErrorGetFile, err)
	assert.Error(t, err)
	filesRepo.AssertExpectations(t)

}

func Test_CreateQueriesAiFail(t *testing.T) {
	var (
		fileId  = 1
		userId  = 1
		prompt  = "Привет,я функция теста системы запросов!"
		mockErr = errors.New("mock error")
	)
	filesRepo := new(mocks.MockFilesRepo)
	queriesRepo := new(mocks.MockQueriesRepo)
	ai := new(mocks.AiMock)
	svc := NewQueriesService(filesRepo, queriesRepo, ai)

	ai.On("Analyze", mock.Anything, mock.Anything).Return("", mockErr)
	filesRepo.On("GetFileById", mock.Anything).Return(models.Files{Path: "testdata/test.csv"}, nil)

	_, err := svc.CreateQueries(fileId, userId, prompt)

	assert.Equal(t, errs.RequestFailed, err)
	assert.Error(t, err)
	ai.AssertExpectations(t)
	filesRepo.AssertExpectations(t)
}

func Test_CreateQueriesFail(t *testing.T) {
	var (
		fileId  = 1
		userId  = 1
		prompt  = "Привет,я функция теста системы запросов!"
		mockErr = errors.New("mock error")
	)
	filesRepo := new(mocks.MockFilesRepo)
	queriesRepo := new(mocks.MockQueriesRepo)
	ai := new(mocks.AiMock)
	svc := NewQueriesService(filesRepo, queriesRepo, ai)
	q := models.Queries{FileId: fileId, UserId: userId, Question: prompt}

	queriesRepo.On("CreateQueries", mock.Anything).Return(q, mockErr)
	filesRepo.On("GetFileById", mock.Anything).Return(models.Files{Path: "testdata/test.csv"}, nil)
	ai.On("Analyze", mock.Anything, mock.Anything).Return("", nil)

	_, err := svc.CreateQueries(fileId, userId, prompt)

	assert.Equal(t, errs.CreateQueryError, err)
	assert.Error(t, err)
	filesRepo.AssertExpectations(t)
	queriesRepo.AssertExpectations(t)
	ai.AssertExpectations(t)
}

func TestGetQueriesByUserIdSuccess(t *testing.T) {
	id := 1
	queriesRepo := new(mocks.MockQueriesRepo)
	filesRepo := new(mocks.MockFilesRepo)
	ai := new(mocks.AiMock)
	svc := NewQueriesService(filesRepo, queriesRepo, ai)

	queriesRepo.On("GetQueriesByUserId", id).Return([]models.Queries{{UserId: id}}, nil)

	query, err := svc.GetQueriesByUserId(id)

	assert.Equal(t, id, query[0].UserId)
	assert.NoError(t, err)
	queriesRepo.AssertExpectations(t)
}

func TestGetQueriesByUserIdFail(t *testing.T) {
	mockErr := errors.New("some error")
	id := 1
	queriesRepo := new(mocks.MockQueriesRepo)
	filesRepo := new(mocks.MockFilesRepo)
	ai := new(mocks.AiMock)
	svc := NewQueriesService(filesRepo, queriesRepo, ai)

	queriesRepo.On("GetQueriesByUserId", id).Return([]models.Queries{{UserId: id}}, mockErr)

	_, err := svc.GetQueriesByUserId(id)

	assert.Equal(t, mockErr, err)
	assert.Error(t, err)
	queriesRepo.AssertExpectations(t)
}
