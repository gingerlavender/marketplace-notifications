package wb

import (
	"fmt"
	"marketplace-notifications/internal/utils/format"
	"strings"
	"time"
)

type Question struct {
	Id             string         `json:"id"`
	Text           string         `json:"text"`
	ProductDetails ProductDetails `json:"productDetails"`
	CreatedDate    time.Time      `json:"createdDate"`
}

func (question Question) FormatMarkdown() string {
	var message strings.Builder

	message.WriteString(fmt.Sprintf("📦  *Товар \\(артикул: %d\\):* %s\n\n", question.ProductDetails.Article, format.EscapeMarkdown(question.ProductDetails.Name)))

	message.WriteString(fmt.Sprintf("💬  *Текст вопроса:* %s\n\n", format.EscapeMarkdown(question.Text)))

	message.WriteString(fmt.Sprintf("🆔  *ID вопроса:* %s\n", format.EscapeMarkdown(question.Id)))
	message.WriteString(fmt.Sprintf("⌚  *Время создания:* %s\n", format.EscapeMarkdown(question.CreatedDate.Format(time.DateTime))))

	return message.String()
}
