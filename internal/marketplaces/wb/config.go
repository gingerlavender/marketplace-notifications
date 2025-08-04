package wb

import "marketplace-notifications/internal/marketplaces"

func GetConfig(JWT string, maxNewQuestions, maxNewFeedbacks int) marketplaces.MarketplaceConfig {
	return marketplaces.MarketplaceConfig{
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
