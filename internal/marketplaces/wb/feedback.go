package wb

import (
	"fmt"
	"marketplace-notifications/internal/utils/format"
	"strings"
	"time"
)

type Feedback struct {
	Id             string         `json:"id"`
	NumberOfStars  int            `json:"productValuation"`
	Pros           string         `json:"pros"`
	Cons           string         `json:"cons"`
	Text           string         `json:"text"`
	ProductDetails ProductDetails `json:"productDetails"`
	CreatedDate    time.Time      `json:"createdDate"`
}

func (feedback Feedback) FormatMarkdown() string {
	var message strings.Builder

	message.WriteString(fmt.Sprintf("📦  *Товар \\(артикул: %d\\):* %s\n\n", feedback.ProductDetails.Article, format.EscapeMarkdown(feedback.ProductDetails.Name)))

	message.WriteString(fmt.Sprintf("📝  *Количество звёзд:* %s\n\n", strings.Repeat("⭐", feedback.NumberOfStars)))

	message.WriteString(fmt.Sprintf("👍  *Достоинства:* %s\n", format.EscapeMarkdown(feedback.Pros)))
	message.WriteString(fmt.Sprintf("👎  *Недостатки:* %s\n", format.EscapeMarkdown(feedback.Cons)))
	message.WriteString(fmt.Sprintf("💬  *Текст вопроса:* %s\n\n", format.EscapeMarkdown(feedback.Text)))

	message.WriteString(fmt.Sprintf("🆔  *ID отзыва:* %s\n", format.EscapeMarkdown(feedback.Id)))
	message.WriteString(fmt.Sprintf("⌚  *Время создания:* %s\n", format.EscapeMarkdown(feedback.CreatedDate.Format(time.DateTime))))

	return message.String()
}
