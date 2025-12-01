package telegram

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"marketplace-notifications/internal/config"
	"marketplace-notifications/internal/marketplaces"
	"marketplace-notifications/internal/marketplaces/wb"
	"marketplace-notifications/internal/marketplaces/yandex"
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
	ChatId    string `json:"chat_id"`
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

func (notifier *TelegramNotifier) SendSummaryNotificationToAllChats(questionsNumber int, feedbacksNumber int) error {
	return notifier.sendNotificationToAllChats(notifier.formatSummaryNotificationMessage(questionsNumber, feedbacksNumber))
}

func (notifier *TelegramNotifier) SendWBQuestionNotificationToAllChats(question wb.Question) error {
	return notifier.sendNotificationToAllChats(notifier.formatUserReactionNotificationMessage(question, marketplaces.Question, "WB"))
}

func (notifier *TelegramNotifier) SendWBFeedbackNotificationToAllChats(feedback wb.Feedback) error {
	return notifier.sendNotificationToAllChats(notifier.formatUserReactionNotificationMessage(feedback, marketplaces.Feedback, "WB"))
}

func (notifier *TelegramNotifier) SendYandexFeedbackNotificationToAllChats(feedback yandex.Feedback) error {
	return notifier.sendNotificationToAllChats(notifier.formatUserReactionNotificationMessage(feedback, marketplaces.Feedback, "Yandex"))
}

func (notifier *TelegramNotifier) sendNotificationToAllChats(text string) error {
	var lastErr error
	var successCount int

	for _, chatId := range notifier.config.ChatIds {
		message := TelegramMessage{
			ChatId:    chatId,
			Text:      text,
			ParseMode: "MarkdownV2",
		}

		if err := notifier.sendMessage(message); err != nil {
			lastErr = err
			log.Printf("[ERROR] Failed to send notification to chat: %s", chatId)
		} else {
			successCount++
		}
	}

	if successCount == 0 && lastErr != nil {
		return fmt.Errorf("Failed to send to all chats. Last error: %w", lastErr)
	}

	return nil
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

func (notifier *TelegramNotifier) formatUserReactionNotificationMessage(userReaction MardownFormatter, reactionType marketplaces.UserReactionType, serviceName string) string {
	var message strings.Builder

	if reactionType == marketplaces.Question {
		message.WriteString(fmt.Sprintf("*❔ Неотвеченный вопрос на %s:*\n\n", serviceName))
	} else {
		message.WriteString(fmt.Sprintf("*💬 Неотвеченный отзыв на %s:*\n\n", serviceName))
	}

	message.WriteString(userReaction.FormatMarkdown())

	return message.String()
}

func (notifier *TelegramNotifier) sendMessageURL() string {
	return fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage", notifier.config.BotToken)
}
