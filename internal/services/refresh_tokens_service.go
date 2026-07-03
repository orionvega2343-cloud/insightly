package services

import (
	"insightly/internal/errs"
	"insightly/internal/models"
	"insightly/internal/repositories"
	"time"
)

type RefreshTokensService interface {
	CreateToken(t models.RefreshToken) (*models.RefreshToken, error)
	GetTokenByValue(token string) (*models.RefreshToken, error)
	DeleteToken(id int) error
}

type RefreshTokensServiceImpl struct {
	Rtr repositories.RefreshTokensRepository
}

func NewRefreshTokensServiceImpl(rtr repositories.RefreshTokensRepository) *RefreshTokensServiceImpl {
	return &RefreshTokensServiceImpl{Rtr: rtr}
}

func (s *RefreshTokensServiceImpl) CreateToken(t models.RefreshToken) (*models.RefreshToken, error) {
	token, err := s.Rtr.CreateToken(t)
	if err != nil {
		return nil, errs.RefreshTokenNotFound
	}
	return token, nil
}

func (s *RefreshTokensServiceImpl) GetTokenByValue(token string) (*models.RefreshToken, error) {
	t, err := s.Rtr.GetTokenByValue(token)
	if err != nil {
		return nil, errs.InvalidGettingToken
	}

	//Проверка, что токен не истек
	if time.Now().After(t.ExpiresAt) {
		return nil, errs.TokenExpiredError
	}

	return t, nil
}

func (s *RefreshTokensServiceImpl) DeleteToken(id int) error {
	return s.Rtr.DeleteToken(id)
}
