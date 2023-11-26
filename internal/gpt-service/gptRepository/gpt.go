package gptrepository

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"regexp"
	"strings"

	"github.com/gocolly/colly"
	"github.com/rs/zerolog/log"
	"github.com/search-platform/gpt-service/internal/gpt-service/models"

	openai "github.com/sashabaranov/go-openai"
)

// BankWebsite представляет данные о веб-сайте банка
type BankWebsite struct {
	Title string `json:"title"`
	Link  string `json:"link"`
}

type GptRepo struct {
	ApiKey string
	client *openai.Client
}

func NewGptRepo(apiKey string) *GptRepo {
	client := openai.NewClient(apiKey)
	return &GptRepo{
		ApiKey: apiKey,
		client: client,
	}
}

// Prompt отправляет запрос в OpenAI и возвращает ответ
func (gpt *GptRepo) GetBankInfo(ctx context.Context, bankName, country string) (*models.BankDetails, error) {

	sites, err := gpt.SearchBankWebsites(ctx, bankName, country)
	if err != nil {
		return nil, err
	}

	sites = sites[:len(sites)/2]

	var banks []models.BankDetails

	for _, site := range sites {
		content, err := gpt.ScrapePageContent(ctx, site.Link)
		if err != nil {
			log.Error().AnErr("search error", err)
			continue
		}
		instruction := "Найди контактные данные банка на странице, представь их в JSON формате. Мне нужны контактные данные только конкретной страны: " + country
		instruction += "{ \"url\": \"\", \"name\": \"\", \"country\": \"\", \"logo_link\": \"\", \"favicon_link\": \"\", \"address\": \"\", \"contacts\": [{\"type\": \"\", \"description\": \"\", \"value\": \"\"}] }"
		instruction += "Type может быть только PHONE (поставить 0) или EMAIL (поставить 1). Информация с сайта ниже: " + content

		// Ограничение длины запроса
		maxLength := 4096
		if len(instruction) > maxLength {
			instruction = instruction[:maxLength]
		}

		resp, err := gpt.client.CreateChatCompletion(
			ctx,
			openai.ChatCompletionRequest{
				Model: openai.GPT3Dot5Turbo,
				Messages: []openai.ChatCompletionMessage{
					{
						Role:    openai.ChatMessageRoleUser,
						Content: instruction,
					},
				},
			},
		)
		if err != nil {
			return nil, err
		}
		fmt.Println(resp.Choices[0].Message.Content)
		bankDetails := models.BankDetails{}
		err = json.Unmarshal([]byte(resp.Choices[0].Message.Content), &bankDetails)
		if err != nil {
			fmt.Println(err)
		}
		banks = append(banks, bankDetails)
	}

	if len(banks) > 0 {
		return &banks[0], nil
	} else {
		return nil, errors.New("no banks found")
	}
}

func (gpt *GptRepo) SearchBankWebsites(ctx context.Context, bankName, country string) ([]BankWebsite, error) {
	apiKey := "AIzaSyA2Lsg8gMg9lBCQHUlT8qFO35LQaai3OLg"
	cx := "a3b9e97770c424185"
	query := fmt.Sprintf("official site bank %s %s contacts", bankName, country)

	// Создание URL запроса
	searchURL := fmt.Sprintf("https://www.googleapis.com/customsearch/v1?q=%s&cx=%s&key=%s", url.QueryEscape(query), cx, apiKey)

	// Выполнение HTTP запроса
	resp, err := http.Get(searchURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Чтение и анализ ответа
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var result struct {
		Items []struct {
			Title string `json:"title"`
			Link  string `json:"link"`
		} `json:"items"`
	}

	err = json.Unmarshal(body, &result)
	if err != nil {
		return nil, err
	}

	// Преобразование результатов в BankWebsite
	var websites []BankWebsite
	for _, item := range result.Items {
		websites = append(websites, BankWebsite{
			Title: item.Title,
			Link:  item.Link,
		})
	}

	return websites, nil
}

// ScrapePageContent посещает указанный URL и извлекает его содержимое
func (gpt *GptRepo) ScrapePageContent(ctx context.Context, url string) (string, error) {
	// Создаем новый экземпляр коллектора
	c := colly.NewCollector()

	var contentBuilder strings.Builder

	// Фильтрация ненужных элементов
	c.OnHTML("h1, h2, h3, h4, h5, h6", func(e *colly.HTMLElement) {
		contentBuilder.WriteString(e.Text + " ")
	})

	// Устанавливаем обработчик для элементов HTML
	// Здесь можно настроить селекторы для конкретных элементов, если нужно
	c.OnHTML("body", func(e *colly.HTMLElement) {
		contentBuilder.WriteString(e.Text)
	})

	// Обработка ошибок
	c.OnError(func(r *colly.Response, err error) {
		fmt.Println("Request URL:", r.Request.URL, "failed with response", "\nError:", err)
	})

	// Посещаем URL
	err := c.Visit(url)
	if err != nil {
		return "", err
	}

	processedContent := preprocessText(contentBuilder.String())

	// Возвращаем всё содержимое страницы в виде строки
	return processedContent, nil
}

// func (gpt *GptRepo) GoogleQuery(ctx context.Context, req string) (string, error) {
// 	instruction := ""

// }

// func (gpt *GptRepo) ParseWebsite(ctx context.Context, url string) (string, error) {

// }

func preprocessText(text string) string {
	// Удаление лишних пробелов и символов переноса строки
	cleanText := regexp.MustCompile(`\s+`).ReplaceAllString(text, " ")
	return strings.TrimSpace(cleanText)
}
