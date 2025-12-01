package yandex

import (
	"fmt"
	"marketplace-notifications/internal/utils/format"
	"strings"
	"time"
)

type Feedback struct {
	Description struct {
		Pros string `json:"advantages"`
		Cons string `json:"disadvantages"`
		Text string `json:"comment"`
	} `json:"description"`
	Statistics struct {
		NumberOfStars int  `json:"rating"`
		Recommended   bool `json:"recommended"`
	} `json:"statistics"`
	Identifiers struct {
		OrderId int `json:"orderId"`
	} `json:"identifiers"`
	Id          int       `json:"feedbackId"`
	CreatedDate time.Time `json:"createdAt"`
}

func (feedback Feedback) FormatMarkdown() string {
	var message strings.Builder

	message.WriteString(fmt.Sprintf("📦  *Заказ с id: %d*\n\n", feedback.Identifiers.OrderId))

	message.WriteString(fmt.Sprintf("📝  *Количество звёзд:* %s\n\n", strings.Repeat("⭐", feedback.Statistics.NumberOfStars)))

	message.WriteString(fmt.Sprintf("👍  *Достоинства:* %s\n", format.EscapeMarkdown(feedback.Description.Pros)))
	message.WriteString(fmt.Sprintf("👎  *Недостатки:* %s\n", format.EscapeMarkdown(feedback.Description.Cons)))
	message.WriteString(fmt.Sprintf("💬  *Текст вопроса:* %s\n\n", format.EscapeMarkdown(feedback.Description.Text)))

	message.WriteString(fmt.Sprintf("🆔  *ID отзыва:* %d\n", feedback.Id))
	message.WriteString(fmt.Sprintf("⌚  *Время создания:* %s\n", format.EscapeMarkdown(feedback.CreatedDate.Format(time.DateTime))))

	return message.String()
}
