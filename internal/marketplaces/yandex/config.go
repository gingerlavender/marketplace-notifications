package yandex

import (
	"fmt"
	"strings"
)

type Config struct {
	APIToken string
	BaseURL  string
	RPS      int
	Burst    int
}

func (config Config) FeedbacksURL(businessId int) string {
	var url strings.Builder

	url.WriteString(config.BaseURL)
	url.WriteString(fmt.Sprintf("/businesses/%d/goods-feedback", businessId))

	return url.String()
}

func GetConfig(APIToken string) Config {
	return Config{
		APIToken: APIToken,
		RPS:      3,
		Burst:    6,
		BaseURL:  "https://api.partner.market.yandex.ru/",
	}
}
