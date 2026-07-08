package ai

import (
	"context"
	"fmt"
	"insightly/internal/errs"

	"github.com/openai/openai-go/v3"
	"github.com/openai/openai-go/v3/responses"
)

type AiAnalyzer interface {
	Analyze(csvData string, prompt string) (string, error)
}

type OpenAIClient struct {
	Client *openai.Client
}

func NewOpenAIClient(client *openai.Client) *OpenAIClient {
	return &OpenAIClient{Client: client}
}

func (a *OpenAIClient) Analyze(csvData string, prompt string) (string, error) {
	//Инициализируем API ключ
	ctx := context.Background()
	client := a.Client

	if csvData == "" || prompt == "" {
		return "", errs.ValueError
	}

	//Отправляем наш запрос, выбираем модель
	response, err := client.Responses.New(ctx, responses.ResponseNewParams{
		Input: responses.ResponseNewParamsInputUnion{OfString: openai.String(fmt.Sprintf("%s,%s", csvData, prompt))},
		Model: openai.ChatModelGPT4,
	})
	if err != nil {
		return "", errs.RequestFailed
	}

	return response.OutputText(), nil
}
