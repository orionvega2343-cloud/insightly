package repositories

import (
	"insightly/internal/models"

	"github.com/jmoiron/sqlx"
)

type UsersRepository interface {
	CreateUser(u models.User) (models.User, error)
	GetUserByEmail(email string) (models.User, error)
}

type UserRepo struct {
	db *sqlx.DB
}

func NewUserRepo(db *sqlx.DB) *UserRepo {
	return &UserRepo{db: db}
}

func (r *UserRepo) CreateUser(u models.User) (models.User, error) {
	err := r.db.QueryRow(`INSERT INTO users(email,password,created_at) VALUES ($1,$2,$3) RETURNING id`, u.Email, u.Password, u.CreatedAt).Scan(&u.Id)
	if err != nil {
		return models.User{}, err
	}
	return u, nil
}

func (r *UserRepo) GetUserByEmail(email string) (models.User, error) {
	var u models.User
	err := r.db.QueryRow(`SELECT id,email,password,created_at FROM users WHERE email=$1`, email).Scan(&u.Id, &u.Email, &u.Password, &u.CreatedAt)
	if err != nil {
		return models.User{}, err
	}
	return u, nil
}
