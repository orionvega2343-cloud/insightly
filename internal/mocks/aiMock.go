package mocks

import "github.com/stretchr/testify/mock"

type AiMock struct {
	mock.Mock
}

func (a *AiMock) Analyze(csvData string, prompt string) (string, error) {
	args := a.Called(csvData, prompt)
	return args.String(0), args.Error(1)
}
