package services

import (
	"crypto/rand"
	"encoding/hex"
	"insightly/internal/errs"
	"insightly/internal/models"
	"insightly/internal/repositories"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

//Not used AI this code

type UserService interface {
	Register(u models.User) (models.User, error)
	Login(email, password, secret string) (string, string, error)
	GenerateAccessToken(userId int, email string) (string, error)
}

type UserServiceImpl struct {
	Ur     repositories.UsersRepository
	Rtr    repositories.RefreshTokensRepository
	Secret string
}

type Claims struct {
	UserId int
	Email  string
	jwt.RegisteredClaims
}

func NewUserServiceImpl(ur repositories.UsersRepository, rtr repositories.RefreshTokensRepository, secret string) *UserServiceImpl {
	return &UserServiceImpl{Ur: ur, Rtr: rtr, Secret: secret}
}

func (s *UserServiceImpl) Register(u models.User) (models.User, error) {
	if u.Email == "" {
		return models.User{}, errs.EmailError
	}

	if u.Password == "" {
		return models.User{}, errs.PasswordError
	}
	//Хэшируем пароль
	hashed, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return models.User{}, errs.InvalidHash
	}
	u.Password = string(hashed)

	user, err := s.Ur.CreateUser(u)
	if err != nil {
		return models.User{}, errs.CreateUserError
	}
	return user, nil
}

func (s *UserServiceImpl) Login(email, password, secret string) (string, string, error) {
	user, err := s.Ur.GetUserByEmail(email)
	if err != nil {
		return "", "", errs.GetEmailError
	}

	//Проверка пароля
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return "", "", errs.InvalidPasswordError
	}

	var c Claims
	c.Email = email
	c.UserId = user.Id
	c.ExpiresAt = jwt.NewNumericDate(time.Now().Add(24 * time.Hour))

	//Генерация JWT токена
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, c)
	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", "", errs.InvalidJWTError
	}

	rToken, err := s.generateRefreshToken()
	if err != nil {
		return "", "", errs.InvalidGenerateToken
	}
	//Сохранение refresh-токена
	refToken := models.RefreshToken{
		UserId:    user.Id,
		Token:     rToken,
		ExpiresAt: time.Now().Add(7 * 24 * time.Hour),
		CreatedAt: time.Now(),
	}

	save, err := s.Rtr.CreateToken(refToken)
	if err != nil {
		return "", "", errs.InvalidSaveToken
	}

	return tokenString, save.Token, nil

}

// Refresh token
func (s *UserServiceImpl) generateRefreshToken() (string, error) {
	bytes := make([]byte, 32)
	_, err := rand.Read(bytes)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

func (s *UserServiceImpl) GenerateAccessToken(userId int, email string) (string, error) {
	var c Claims
	c.UserId = userId
	c.Email = email
	c.ExpiresAt = jwt.NewNumericDate(time.Now().Add(24 * time.Hour))
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, c)
	tokenString, err := token.SignedString([]byte(s.Secret))
	if err != nil {
		return "", errs.InvalidJWTError
	}
	return tokenString, nil
}
