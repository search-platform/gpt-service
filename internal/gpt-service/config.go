package gptservice

type ServiceConfig struct {
	ApiKey             string `env:"GPT_4_API_KEY"`
	CustomSearchApiKey string `env:"GOOGLE_API_KEY"`
	CustomSearchCX     string `env:"GOOGLE_CUSTOM_SEARCH_ENGINE_ID"`
}
