package services

import (
	"errors"
	"insightly/internal/errs"
	"insightly/internal/mocks"
	"insightly/internal/models"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func Test_CreateTokenSuccess(t *testing.T) {
	//Arrange
	tokensRepo := new(mocks.MockTokensRepo)
	svc := NewRefreshTokensServiceImpl(tokensRepo)
	m := models.RefreshToken{Token: "abc123"}

	tokensRepo.On("CreateToken", mock.Anything).Return(&m, nil)

	//Act
	token, err := svc.CreateToken(m)

	//Assert
	assert.Equal(t, m.Token, token.Token)
	assert.NoError(t, err)
	tokensRepo.AssertExpectations(t)
}

func Test_CreateTokenFail(t *testing.T) {
	//Arrange
	tokensRepo := new(mocks.MockTokensRepo)
	svc := NewRefreshTokensServiceImpl(tokensRepo)
	m := models.RefreshToken{Token: "abc123"}

	tokensRepo.On("CreateToken", mock.Anything).Return(&m, errors.New("some error"))

	//Act
	_, err := svc.CreateToken(m)

	//Assert
	assert.Equal(t, errs.RefreshTokenNotFound, err)
	assert.Error(t, err)
	tokensRepo.AssertExpectations(t)
}

func Test_GetTokenByValueSuccess(t *testing.T) {
	//Arrange
	tokensRepo := new(mocks.MockTokensRepo)
	svc := NewRefreshTokensServiceImpl(tokensRepo)
	m := models.RefreshToken{Token: "abc123", ExpiresAt: time.Now().Add(1 * time.Hour)}

	tokensRepo.On("GetTokenByValue", mock.Anything).Return(&m, nil)

	//Act
	token, err := svc.GetTokenByValue(m.Token)

	assert.Equal(t, m.Token, token.Token)
	assert.NoError(t, err)
	tokensRepo.AssertExpectations(t)
}

func Test_GetTokenByValueFail(t *testing.T) {
	//Arrange
	tokensRepo := new(mocks.MockTokensRepo)
	svc := NewRefreshTokensServiceImpl(tokensRepo)
	m := models.RefreshToken{Token: "abc123", ExpiresAt: time.Now().Add(1 * time.Hour)}

	tokensRepo.On("GetTokenByValue", mock.Anything).Return(&m, errors.New("some error"))

	//Act
	_, err := svc.GetTokenByValue(m.Token)

	assert.Equal(t, errs.InvalidGettingToken, err)
	assert.Error(t, err)
	tokensRepo.AssertExpectations(t)
}

func Test_GetTokenByValueExpired(t *testing.T) {
	//Arrange
	tokensRepo := new(mocks.MockTokensRepo)
	svc := NewRefreshTokensServiceImpl(tokensRepo)
	m := models.RefreshToken{Token: "abc123", ExpiresAt: time.Now().Add(-1 * time.Hour)}

	tokensRepo.On("GetTokenByValue", mock.Anything).Return(&m, nil)

	//Act
	_, err := svc.GetTokenByValue(m.Token)

	assert.Equal(t, errs.TokenExpiredError, err)
	assert.Error(t, err)
	tokensRepo.AssertExpectations(t)
}
