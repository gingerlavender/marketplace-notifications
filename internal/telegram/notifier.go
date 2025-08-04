package telegram

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"marketplace-notifications/internal/config"
	"marketplace-notifications/internal/marketplaces"
	"marketplace-notifications/internal/marketplaces/wb"
	"net/http"
	"strings"

	"golang.org/x/time/rate"
)

type MardownFormatter interface {
	FormatMarkdown() string
}

type TelegramNotifier struct {
	config          *config.TelegramConfig
	httpClient      *http.Client
	telegramLimiter *rate.Limiter
}

type TelegramMessage struct {
	ChatID    string `json:"chat_id"`
	Text      string `json:"text"`
	ParseMode string `json:"parse_mode"`
}

func NewTelegramNotifier(config *config.TelegramConfig) *TelegramNotifier {
	telegramLimiter := rate.NewLimiter(rate.Limit(config.RPS), config.RPS)

	return &TelegramNotifier{
		config: config,
		httpClient: &http.Client{
			Timeout: config.Timeout,
		},
		telegramLimiter: telegramLimiter,
	}
}

func (notifier *TelegramNotifier) SendSummaryNotification(questionsNumber, feedbacksNumber int) error {
	message := TelegramMessage{
		ChatID:    notifier.config.ChatId,
		Text:      notifier.formatSummaryNotificationMessage(questionsNumber, feedbacksNumber),
		ParseMode: "MarkdownV2",
	}

	return notifier.sendMessage(message)
}

func (notifier *TelegramNotifier) SendQuestionNotification(question wb.Question) error {
	message := TelegramMessage{
		ChatID:    notifier.config.ChatId,
		Text:      notifier.formatUserReactionNotificationMessage(question, marketplaces.Question),
		ParseMode: "MarkdownV2",
	}

	return notifier.sendMessage(message)
}

func (notifier *TelegramNotifier) SendFeedbackNotification(feedback wb.Feedback) error {
	message := TelegramMessage{
		ChatID:    notifier.config.ChatId,
		Text:      notifier.formatUserReactionNotificationMessage(feedback, marketplaces.Feedback),
		ParseMode: "MarkdownV2",
	}

	return notifier.sendMessage(message)
}

func (notifier *TelegramNotifier) sendMessage(message TelegramMessage) error {
	if err := notifier.telegramLimiter.Wait(context.Background()); err != nil {
		return fmt.Errorf("Telegram rate limiter error: %w", err)
	}

	jsonData, err := json.Marshal(message)
	if err != nil {
		return fmt.Errorf("failed to marshal message: %w", err)
	}

	resp, err := notifier.httpClient.Post(notifier.sendMessageURL(), "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("failed to send message: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("Telegram returned status %d instead of 200: %s", resp.StatusCode, body)
	}

	return nil
}

func (notifier *TelegramNotifier) formatSummaryNotificationMessage(questionsNumber, feedbacksNumber int) string {
	var message strings.Builder

	message.WriteString(fmt.Sprintf("🔔 *Пользователи ждут вашего ответа\\!* 🔔\n\n"))

	message.WriteString(fmt.Sprintf("*🗓️ На данный момент у вас:*\n\n"))

	message.WriteString(fmt.Sprintf("❔ Неотвеченных *вопросов*: %d\n", questionsNumber))
	message.WriteString(fmt.Sprintf("💬 Неотвеченных *отзывов*: %d\n\n", feedbacksNumber))

	message.WriteString(fmt.Sprintf("📃 Полный список в сообщениях ниже:\n"))

	return message.String()
}

func (notifier *TelegramNotifier) formatUserReactionNotificationMessage(userReaction MardownFormatter, reactionType marketplaces.UserReactionType) string {
	var message strings.Builder

	if reactionType == marketplaces.Question {
		message.WriteString(fmt.Sprintf("*❔ Неотвеченный вопрос:*\n\n"))
	} else {
		message.WriteString(fmt.Sprintf("*💬 Неотвеченный отзыв:*\n\n"))
	}

	message.WriteString(userReaction.FormatMarkdown())

	return message.String()
}

func (notifier *TelegramNotifier) sendMessageURL() string {
	return fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage", notifier.config.BotToken)
}
