package errs

import "errors"

// Ошибки сервисов
var (
	EmailError           = errors.New("invalid email")
	PasswordError        = errors.New("invalid password")
	InvalidHash          = errors.New("invalid hash")
	CreateUserError      = errors.New("user not found")
	GetEmailError        = errors.New("error getting mail")
	InvalidPasswordError = errors.New("invalid password")
	InvalidJWTError      = errors.New("invalid token")
	InvalidRefreshtoken  = errors.New("invalid refresh token")
	InvalidGenerateToken = errors.New("invalid token")
	InvalidSaveToken     = errors.New("invalid save token")
	RefreshTokenNotFound = errors.New("token not found")
	InvalidGettingToken  = errors.New("invalid token")
	TokenExpiredError    = errors.New("excessive waiting time")
	FailedSaveFile       = errors.New("failed to save file")
	FailedCreateFile     = errors.New("failed to create file")
)

// Ошибки API
var (
	RequestFailed = errors.New("request failed")
	ValueError    = errors.New("strings csvData and/or prompt can't be empty")
)

// Ошибка Queries
var (
	ErrorGetFile     = errors.New("error getting file")
	CreateQueryError = errors.New("error creating query")
)

// Ошибки парсера
var (
	ErrorOpenFile  = errors.New("error opening file")
	ErrorReadFile  = errors.New("error reading file")
	InvalidFormat  = errors.New("invalid format")
	ErrorEmptyFile = errors.New("empty file")
)
