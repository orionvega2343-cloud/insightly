package services

import (
	"insightly/internal/ai"
	"insightly/internal/errs"
	"insightly/internal/models"
	"insightly/internal/parser"
	"insightly/internal/repositories"
)

type QueriesService interface {
	CreateQueries(fileId, userId int, prompt string) (models.Queries, error)
	GetQueriesByUserId(userid int) ([]models.Queries, error)
}

type QueriesServiceImpl struct {
	F  repositories.FilesRepository
	Q  repositories.QueriesRepository
	Ai ai.AiAnalyzer
}

func NewQueriesService(f repositories.FilesRepository, q repositories.QueriesRepository, ai ai.AiAnalyzer) *QueriesServiceImpl {
	return &QueriesServiceImpl{F: f, Q: q, Ai: ai}
}

func (s *QueriesServiceImpl) CreateQueries(fileId, userId int, prompt string) (models.Queries, error) {
	// Получаем файл по fileId
	file, err := s.F.GetFileById(fileId)
	if err != nil {
		return models.Queries{}, errs.ErrorGetFile
	}

	csvData, err := parser.CsvParser(file.Path)
	if err != nil {
		return models.Queries{}, errs.InvalidPath
	}

	// Вызов AI
	answer, err := s.Ai.Analyze(csvData, prompt)
	if err != nil {
		return models.Queries{}, errs.RequestFailed
	}

	q := models.Queries{
		UserId:   userId,
		FileId:   fileId,
		Question: prompt,
		Answer:   answer,
	}

	result, err := s.Q.CreateQueries(q)
	if err != nil {
		return models.Queries{}, errs.CreateQueryError
	}

	return result, nil
}

func (s *QueriesServiceImpl) GetQueriesByUserId(userId int) ([]models.Queries, error) {
	return s.Q.GetQueriesByUserId(userId)
}
