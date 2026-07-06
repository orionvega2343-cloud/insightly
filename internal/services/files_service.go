package services

import (
	"fmt"
	"insightly/internal/errs"
	"insightly/internal/models"
	"insightly/internal/repositories"
	"os"
	"path/filepath"
	"strings"

	"github.com/google/uuid"
)

type FilesService interface {
	CreateFiles(userId int, filename string, data []byte) (models.Files, error)
	GetFilesByUserId(userId int) ([]models.Files, error)
	GetFileById(id int) (models.Files, error)
	DeleteFile(id int) error
}

type FilesServiceImpl struct {
	F repositories.FilesRepository
}

func NewFilesService(f repositories.FilesRepository) *FilesServiceImpl {
	return &FilesServiceImpl{F: f}
}

func (s *FilesServiceImpl) CreateFiles(userId int, filename string, data []byte) (models.Files, error) {
	//Проверяем расширение файла,
	//принимаем только .csv
	ext := filepath.Ext(filename)
	if !strings.EqualFold(ext, ".csv") {
		return models.Files{}, errs.InvalidFileExtension
	}

	id := uuid.New().String()

	//Собираем пути
	dir := fmt.Sprintf("uploads/%d", userId)
	err := os.MkdirAll(dir, 0755)
	if err != nil {
		return models.Files{}, err
	}

	//записываем в байты
	fullPath := filepath.Join(dir, id+ext)
	err = os.WriteFile(fullPath, data, 0666)
	if err != nil {
		return models.Files{}, err
	}

	f := models.Files{
		UserId: userId,
		Name:   filename,
		Path:   fullPath,
	}

	created, err := s.F.CreateFiles(f)
	if err != nil {
		return models.Files{}, errs.FailedCreateFile
	}

	return created, nil
}

func (s *FilesServiceImpl) GetFilesByUserId(userId int) ([]models.Files, error) {
	return s.F.GetFilesByUserId(userId)
}

func (s *FilesServiceImpl) GetFileById(id int) (models.Files, error) {
	return s.F.GetFileById(id)
}

func (s *FilesServiceImpl) DeleteFile(id int) error {
	return s.F.DeleteFile(id)
}
