package marketplaces

import (
	"strings"
)

type MarketplaceConfig struct {
	JWT             string
	RPS             int
	Burst           int
	BaseURL         string
	QuestionsPath   string
	FeedbacksPath   string
	MaxNewQuestions int
	MaxNewFeedbacks int
}

func (config MarketplaceConfig) QuestionsURL() string {
	var url strings.Builder

	url.WriteString(config.BaseURL)
	url.WriteString(config.QuestionsPath)

	return url.String()
}

func (config MarketplaceConfig) FeedbacksURL() string {
	var url strings.Builder

	url.WriteString(config.BaseURL)
	url.WriteString(config.FeedbacksPath)

	return url.String()
}
