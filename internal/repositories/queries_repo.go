package repositories

import (
	"insightly/internal/models"

	"github.com/jmoiron/sqlx"
)

type QueriesRepository interface {
	CreateQueries(q models.Queries) (models.Queries, error)
	GetQueriesByUserId(id int) ([]models.Queries, error)
}

type QueriesRepo struct {
	db *sqlx.DB
}

func NewQueriesRepo(db *sqlx.DB) *QueriesRepo {
	return &QueriesRepo{db: db}
}

func (r *QueriesRepo) CreateQueries(q models.Queries) (models.Queries, error) {
	err := r.db.QueryRow(`INSERT INTO queries (user_id,file_id,created_at,question,answer) VALUES($1,$2,$3,$4,$5) RETURNING id`, q.UserId, q.FileId, q.CreatedAt, q.Question, q.Answer).Scan(&q.Id)
	if err != nil {
		return models.Queries{}, err
	}
	return q, nil
}

func (r *QueriesRepo) GetQueriesByUserId(id int) ([]models.Queries, error) {
	var q []models.Queries
	err := r.db.Select(&q, `SELECT id, user_id, file_id, created_at, question, answer FROM queries WHERE user_id = $1`, id)
	if err != nil {
		return []models.Queries{}, err
	}
	return q, nil
}
