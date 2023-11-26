package gptrepository

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"time"
)

const (
	gpt4APIURL = "https://api.openai.com/v1/engines/gpt-4/completions"
)

type Prompt struct {
	Prompt    string `json:"prompt"`
	MaxTokens int    `json:"max_tokens"`
}

type GptRepo struct {
	ApiKey string
}

func NewGptRepo(apiKey string) *GptRepo {
	return &GptRepo{
		ApiKey: apiKey,
	}
}

func (gpt *GptRepo) Prompt(ctx context.Context, prompt string) (string, error) {
	gpt4Req := Prompt{
		Prompt:    prompt,
		MaxTokens: 150,
	}

	reqBody, err := json.Marshal(gpt4Req)
	if err != nil {
		return "", err
	}

	// Создаем HTTP-запрос
	httpRequest, err := http.NewRequestWithContext(ctx, "POST", gpt4APIURL, bytes.NewBuffer(reqBody))
	if err != nil {
		return "", err
	}
	httpRequest.Header.Set("Content-Type", "application/json")
	httpRequest.Header.Set("Authorization", "Bearer "+gpt.ApiKey)

	// Отправляем запрос
	httpClient := &http.Client{Timeout: 10 * time.Second}
	resp, err := httpClient.Do(httpRequest)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	// Обрабатываем ответ
	var gpt4Resp struct {
		Choices []struct {
			Text string `json:"text"`
		} `json:"choices"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&gpt4Resp); err != nil {
		return "", err
	}

	// Формируем и возвращаем ответ
	return gpt4Resp.Choices[0].Text, nil
}
