package services

import (
	"insightly/internal/errs"
	"insightly/internal/models"
	"insightly/internal/repositories"
	"os"
)

type FilesService interface {
	CreateFiles(f models.Files) (models.Files, error)
	GetFilesByUserId(userId int) ([]models.Files, error)
	DeleteFile(id int) error
}

type FilesServiceImpl struct {
	F repositories.FilesRepository
}

func NewFilesService(f repositories.FilesRepository) *FilesServiceImpl {
	return &FilesServiceImpl{F: f}
}

func (s *FilesServiceImpl) CreateFiles(f models.Files) (models.Files, error) {
	data := []byte("Hello world")
	err := os.WriteFile(f.Path, data, 0666)
	if err != nil {
		return f, errs.FailedSaveFile
	}

	file, err := s.F.CreateFiles(f)
	if err != nil {
		return f, errs.FailedCreateFile
	}
	return file, nil
}

func (s *FilesServiceImpl) GetFilesByUserId(userId int) ([]models.Files, error) {
	return s.F.GetFilesByUserId(userId)
}

func (s *FilesServiceImpl) DeleteFile(id int) error {
	return s.F.DeleteFile(id)
}
