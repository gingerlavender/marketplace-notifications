package wb

import (
	"strings"
)

type Config struct {
	JWT             string
	RPS             int
	Burst           int
	BaseURL         string
	QuestionsPath   string
	FeedbacksPath   string
	MaxNewQuestions int
	MaxNewFeedbacks int
}

func (config Config) QuestionsURL() string {
	var url strings.Builder

	url.WriteString(config.BaseURL)
	url.WriteString(config.QuestionsPath)

	return url.String()
}

func (config Config) FeedbacksURL() string {
	var url strings.Builder

	url.WriteString(config.BaseURL)
	url.WriteString(config.FeedbacksPath)

	return url.String()
}

func GetConfig(JWT string, maxNewQuestions, maxNewFeedbacks int) Config {
	return Config{
		JWT:             JWT,
		RPS:             3,
		Burst:           6,
		BaseURL:         "https://feedbacks-api.wildberries.ru/api/v1/",
		QuestionsPath:   "questions",
		FeedbacksPath:   "feedbacks",
		MaxNewQuestions: maxNewQuestions,
		MaxNewFeedbacks: maxNewFeedbacks,
	}
}
