package repositories

import (
	"insightly/internal/models"

	"github.com/jmoiron/sqlx"
)

type RefreshTokensRepository interface {
	CreateToken(t models.RefreshToken) (*models.RefreshToken, error)
	GetTokenByValue(token string) (*models.RefreshToken, error)
	DeleteToken(id int) error
}

type RefreshTokensRepo struct {
	db *sqlx.DB
}

func NewRefreshTokensRepo(db *sqlx.DB) *RefreshTokensRepo {
	return &RefreshTokensRepo{db: db}
}

func (r *RefreshTokensRepo) CreateToken(t models.RefreshToken) (*models.RefreshToken, error) {
	err := r.db.QueryRow(`INSERT INTO refresh_tokens(user_id,expires_at,token,created_at) VALUES($1,$2,$3,$4) RETURNING id`, t.UserId, t.ExpiresAt, t.Token, t.CreatedAt).Scan(&t.Id)
	if err != nil {
		return nil, err
	}
	return &t, nil
}

func (r *RefreshTokensRepo) GetTokenByValue(token string) (*models.RefreshToken, error) {
	var t models.RefreshToken
	err := r.db.QueryRow(`SELECT id,user_id,expires_at,token,created_at FROM refresh_tokens WHERE token=$1`, token).Scan(&t.Id, &t.UserId, &t.ExpiresAt, &t.Token, &t.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &t, nil
}

func (r *RefreshTokensRepo) DeleteToken(id int) error {
	_, err := r.db.Exec("DELETE FROM refresh_tokens WHERE id=$1", id)
	if err != nil {
		return err
	}
	return nil
}
