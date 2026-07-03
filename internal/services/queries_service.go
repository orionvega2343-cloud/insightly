package services

import (
	"insightly/internal/ai"
	"insightly/internal/errs"
	"insightly/internal/models"
	"insightly/internal/repositories"
)

type QueriesService interface {
	CreateQueries(fileId, userId int, prompt string) (models.Queries, error)
	GetQueriesByUserId(userid int) ([]models.Queries, error)
}

type QueriesServiceImpl struct {
	F  repositories.FilesRepository
	Q  repositories.QueriesRepository
	Ai *ai.OpenAIClient
}

func NewQueriesService(f repositories.FilesRepository, q repositories.QueriesRepository, ai *ai.OpenAIClient) *QueriesServiceImpl {
	return &QueriesServiceImpl{F: f, Q: q, Ai: ai}
}

func (s *QueriesServiceImpl) CreateQueries(fileId, userId int, prompt string) (models.Queries, error) {
	// Получаем файл по fileId
	_, err := s.F.GetFileById(fileId)
	if err != nil {
		return models.Queries{}, errs.ErrorGetFile
	}

	// TODO: вызов парсера CSV по file.Path
	csvData := "заглушка"

	// Вызов AI
	answer, err := s.Ai.Analyze(csvData, prompt)
	if err != nil {
		return models.Queries{}, errs.RequestFailed
	}

	// Формируем модель
	q := models.Queries{
		UserId:   userId,
		FileId:   fileId,
		Question: prompt,
		Answer:   answer,
	}

	// Сохраняем в БД
	result, err := s.Q.CreateQueries(q)
	if err != nil {
		return models.Queries{}, errs.CreateQueryError
	}

	return result, nil
}

func (s *QueriesServiceImpl) GetQueriesByUserId(userId int) ([]models.Queries, error) {
	return s.Q.GetQueriesByUserId(userId)
}
