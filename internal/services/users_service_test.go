package services

import (
	"errors"
	"insightly/internal/mocks"
	"insightly/internal/models"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"golang.org/x/crypto/bcrypt"
)

func Test_CreateUser(t *testing.T) {
	//Arrange
	userRepo := new(mocks.MockUserRepo)
	tokensRepo := new(mocks.MockTokensRepo)
	var secret string
	svc := NewUserServiceImpl(userRepo, tokensRepo, secret)
	u := models.User{Email: "test@test.com", Password: "12345"}

	userRepo.On("CreateUser", mock.Anything).Return(u, nil)

	//Act
	user, err := svc.Register(u)
	if err != nil {
		t.Error(err)
	}

	//Assert
	assert.Equal(t, u.Email, user.Email)
	assert.NoError(t, err)
	userRepo.AssertExpectations(t)
}

func Test_CreateUserFail(t *testing.T) {
	//Arrange
	userRepo := new(mocks.MockUserRepo)
	tokensRepo := new(mocks.MockTokensRepo)
	var secret string
	svc := NewUserServiceImpl(userRepo, tokensRepo, secret)
	u := models.User{Email: "test@test.com", Password: "12345"}

	userRepo.On("CreateUser", mock.Anything).Return(u, errors.New("error"))

	//Act
	_, err := svc.Register(u)

	//Assert
	assert.Error(t, err)
	userRepo.AssertExpectations(t)
}

func Test_LoginSuccess(t *testing.T) {
	//Arrange
	var (
		email    = "test@test.com"
		secret   = "test-secret"
		password = "12345"
	)
	hashed, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	userRepo := new(mocks.MockUserRepo)
	tokensRepo := new(mocks.MockTokensRepo)
	svc := NewUserServiceImpl(userRepo, tokensRepo, secret)

	u := models.User{Email: email, Password: string(hashed)}
	userRepo.On("GetUserByEmail", mock.Anything).Return(u, nil)
	tokensRepo.On("CreateToken", mock.Anything).Return(&models.RefreshToken{Token: "test-refresh-token"}, nil)

	//Act
	accessToken, refreshToken, err := svc.Login(email, password, secret)

	//Assert
	assert.NoError(t, err)
	assert.NotEmpty(t, accessToken)
	assert.NotEmpty(t, refreshToken)
	userRepo.AssertExpectations(t)
	tokensRepo.AssertExpectations(t)
}

func Test_LoginFail(t *testing.T) {
	//Arrange
	var (
		email         = "test@test.com"
		secret        = "test-secret"
		password      = "12345"
		wrongPassword = "wrong"
	)
	hashed, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	userRepo := new(mocks.MockUserRepo)
	tokensRepo := new(mocks.MockTokensRepo)
	svc := NewUserServiceImpl(userRepo, tokensRepo, secret)

	u := models.User{Email: email, Password: string(hashed)}
	userRepo.On("GetUserByEmail", mock.Anything).Return(u, nil)

	//Act
	accessToken, refreshToken, err := svc.Login(email, wrongPassword, secret)

	//Assert
	assert.Error(t, err)
	assert.Empty(t, accessToken)
	assert.Empty(t, refreshToken)
	userRepo.AssertExpectations(t)
}
