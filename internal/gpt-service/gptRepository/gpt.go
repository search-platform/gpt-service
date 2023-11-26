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
	"github.com/k3a/html2text"
	"github.com/rs/zerolog/log"
	"github.com/search-platform/gpt-service/internal/gpt-service/models"

	openai "github.com/sashabaranov/go-openai"
)

// BankWebsite представляет данные о веб-сайте банка
type BankWebsite struct {
	Title string `json:"title"`
	Link  string `json:"link"`
	Icon  string `json:"icon"`
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

	sitesWithPhone, err := gpt.SearchBankWebsites(ctx, bankName, country, "phone")
	if err != nil {
		return nil, err
	}
	if len(sitesWithPhone) > 3 {
		sitesWithPhone = sitesWithPhone[:3]
	}
	var banks []models.BankDetails

	for _, site := range sitesWithPhone {
		content, favicon, err := gpt.ScrapePageContent(ctx, site.Link)
		if err != nil {
			log.Error().AnErr("search error", err)
			continue
		}
		bankJson, err := json.Marshal(site)
		if err != nil {
			return nil, err
		}
		instruction := "Найди контактные данные банка в тексте."
		instruction += "У меня есть некоторые данные для тебя по этому банку: " + string(bankJson)
		instruction += "Вот вероятная ссылка на favicon: " + favicon
		instruction += ". Также, очень важно найти ссылку на логотип компании в html документе"
		instruction += "Ты должен представить их обязательно в JSON формате. Нельзя ничего сообщать, кроме JSON ответа. "
		instruction += "{ \"url\": \"\", \"name\": \"\", \"country\": \"\", \"logo_link\": \"\", \"favicon_link\": \"\", \"address\": \"\", \"contacts\": [{\"type\": \"\", \"description\": \"\", \"value\": \"\"}] }"
		instruction += "Type может быть только PHONE или EMAIL. Я запрещаю возвращать данные, которых нет на странице. Информация с сайта ниже: " + content

		// Ограничение длины запроса
		maxLength := 4096
		if len(instruction) > maxLength {
			instruction = instruction[:maxLength]
		}

		resp, err := gpt.client.CreateChatCompletion(
			ctx,
			openai.ChatCompletionRequest{
				Model: openai.GPT4TurboPreview,
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
		rawData := resp.Choices[0].Message.Content

		prefix := "```json"
		suffix := "```"
		// Уберем префикс
		rawData = strings.TrimPrefix(rawData, prefix)
		// Уберем постфикс
		rawData = strings.TrimSuffix(rawData, suffix)

		fmt.Println(rawData)
		bankDetails := models.BankDetails{}
		err = json.Unmarshal([]byte(rawData), &bankDetails)
		if err != nil {
			fmt.Println(err)
		}
		banks = append(banks, bankDetails)
	}

	jsonBanks, err := json.Marshal(banks)
	if err != nil {
		return nil, err
	}

	instruction := "Я тебе передаю массив с данными банка " + bankName + " в стране " + country
	instruction += "Тебе запрещено писать что-либо, кроме JSON. Твоя задача: сделать один наиболее вероятный объект с данными банка в JSON, используй в том числе собственную базу знаний: "
	instruction += "{ \"url\": \"\", \"name\": \"\", \"country\": \"\", \"logo_link\": \"\", \"favicon_link\": \"\", \"address\": \"\", \"contacts\": [{\"type\": \"\", \"description\": \"\", \"value\": \"\"}] }"
	instruction += " Известные данные банка: " + string(jsonBanks)

	resp, err := gpt.client.CreateChatCompletion(
		ctx,
		openai.ChatCompletionRequest{
			Model: openai.GPT4TurboPreview,
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

	rawData := resp.Choices[0].Message.Content

	prefix := "```json"
	suffix := "```"
	// Уберем префикс
	rawData = strings.TrimPrefix(rawData, prefix)
	// Уберем постфикс
	rawData = strings.TrimSuffix(rawData, suffix)

	fmt.Println(rawData)

	bankDetails := models.BankDetails{}
	err = json.Unmarshal([]byte(rawData), &bankDetails)
	if err != nil {
		fmt.Println(err)
	}
	return &bankDetails, nil
}

func (gpt *GptRepo) SearchBankWebsites(ctx context.Context, bankName, country, target string) ([]BankWebsite, error) {
	apiKey := "AIzaSyA2Lsg8gMg9lBCQHUlT8qFO35LQaai3OLg"
	cx := "a3b9e97770c424185"

	instruction := "Мне нужен JSON с полями `website`: official_bank_domain_name, `contacts`: слова телефон email на языке страны " + country
	instruction += "для банка " + bankName + ", ты должен заменить official_bank_domain_name на доменное имя банка" + bankName
	instruction += " в этой стране в ответном JSON. "
	instruction += " Запрещено писать что то кроме этого JSON в ответе"
	gptResp, err := gpt.client.CreateChatCompletion(
		ctx,
		openai.ChatCompletionRequest{
			Model: openai.GPT4TurboPreview,
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

	var bankQuery struct {
		Website  string `json:"website"`
		Contacts string `json:"contacts"`
	}

	jsonGptResp := gptResp.Choices[0].Message.Content

	prefix := "```json"
	suffix := "```"
	// Уберем префикс
	jsonGptResp = strings.TrimPrefix(jsonGptResp, prefix)
	// Уберем постфикс
	jsonGptResp = strings.TrimSuffix(jsonGptResp, suffix)

	json.Unmarshal([]byte(jsonGptResp), &bankQuery)

	query := fmt.Sprintf("site:%s %s", bankQuery.Website, bankQuery.Contacts)

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

	fmt.Printf("google response: %s", string(body))

	var result struct {
		Items []struct {
			Title   string `json:"title"`
			Link    string `json:"link"`
			PageMap struct {
				CSEThumbnail []struct {
					Src string `json:"src"`
				} `json:"cse_thumbnail"`
			} `json:"pagemap"`
		} `json:"items"`
	}

	err = json.Unmarshal(body, &result)
	if err != nil {
		return nil, err
	}

	// Преобразование результатов в BankWebsite
	var websites []BankWebsite
	for _, item := range result.Items {
		iconURL := ""
		if len(item.PageMap.CSEThumbnail) > 0 {
			iconURL = item.PageMap.CSEThumbnail[0].Src
		}
		websites = append(websites, BankWebsite{
			Title: item.Title,
			Link:  item.Link,
			Icon:  iconURL,
		})
	}

	return websites, nil
}

func (gpt *GptRepo) ScrapePageContent(ctx context.Context, url string) (string, string, error) {
	// Инициализация коллектора Colly
	c := colly.NewCollector(
		colly.Async(true),
	)

	var content, favicon string

	// Обработка ошибок
	c.OnError(func(_ *colly.Response, err error) {
		content = ""
		favicon = ""
	})

	// Обработка HTML для получения контента
	c.OnHTML("body", func(e *colly.HTMLElement) {
		rawHtml := e.DOM.ParentsUntil("~").Text()

		// Обработка и преобразование HTML в текст
		processedHtml := preprocessText(rawHtml)
		processedHtml = strings.ReplaceAll(processedHtml, "<>", "")
		processedHtml, err := removeStyleTags(processedHtml)
		if err != nil {
			content = ""
			return
		}
		content = html2text.HTML2Text(processedHtml)
	})

	// Извлекаем Favicon
	c.OnHTML("link[rel='icon'], link[rel='shortcut icon']", func(e *colly.HTMLElement) {
		favicon = e.Attr("href")
		if !strings.HasPrefix(favicon, "http") {
			favicon = e.Request.AbsoluteURL(favicon)
		}
	})

	// Отправляем запрос
	c.Visit(url)

	// Ожидаем завершения всех асинхронных операций
	c.Wait()

	if content == "" || favicon == "" {
		return "", "", errors.New("Failed to scrape the content or favicon")
	}

	content = preprocessText(content)

	// fmt.Println(content)

	return content, favicon, nil
}

func removeStyleTags(input string) (string, error) {
	// Компилируем регулярное выражение для поиска тегов <style>
	re, err := regexp.Compile(`<style.*?</style>`)
	if err != nil {
		return "", err
	}

	// Заменяем все найденные вхождения на пустую строку
	result := re.ReplaceAllString(input, "")

	re, err = regexp.Compile(`<script.*?</script>`)
	if err != nil {
		return "", err
	}
	result = re.ReplaceAllString(result, "")

	return result, nil
}

func preprocessText(text string) string {
	// Удаление лишних пробелов и символов переноса строки
	compact := strings.ReplaceAll(text, "\n", "")
	compact = strings.ReplaceAll(compact, "\r", "")
	compact = strings.Join(strings.Fields(compact), " ")

	cleanText := regexp.MustCompile(`\s+`).ReplaceAllString(compact, " ")
	return strings.TrimSpace(cleanText)
}
