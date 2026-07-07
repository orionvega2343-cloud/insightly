package repositories

import (
	"insightly/internal/models"

	"github.com/jmoiron/sqlx"
)

type FilesRepository interface {
	CreateFiles(f models.Files) (models.Files, error)
	GetFilesByUserId(userId int) ([]models.Files, error)
	GetFileById(id int) (models.Files, error)
	DeleteFile(id int) error
}

type FilesRepo struct {
	db *sqlx.DB
}

func NewFilesRepo(db *sqlx.DB) *FilesRepo {
	return &FilesRepo{db: db}
}

func (r *FilesRepo) CreateFiles(f models.Files) (models.Files, error) {
	err := r.db.QueryRow(`INSERT INTO files(user_id,filename,filepath) VALUES($1,$2,$3) RETURNING id`, f.UserId, f.Name, f.Path).Scan(&f.Id)
	if err != nil {
		return f, err
	}
	return f, nil
}

func (r *FilesRepo) GetFilesByUserId(UserId int) ([]models.Files, error) {
	var files []models.Files
	err := r.db.Select(&files, `SELECT id,user_id,filename,filepath FROM files WHERE user_id = $1`, UserId)
	if err != nil {
		return []models.Files{}, err
	}
	return files, nil
}

func (r *FilesRepo) DeleteFile(id int) error {
	_, err := r.db.Exec("DELETE FROM files WHERE id = $1", id)
	if err != nil {
		return err
	}
	return nil
}

func (r *FilesRepo) GetFileById(id int) (models.Files, error) {
	var f models.Files
	err := r.db.QueryRow(`SELECT id, user_id, filename, filepath FROM files WHERE id = $1`, id).Scan(&f.Id, &f.UserId, &f.Name, &f.Path)
	if err != nil {
		return models.Files{}, err
	}
	return f, nil
}
