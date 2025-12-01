package yandex

type NotificationBase struct {
	NotificationType string `json:"notificationType"`
}

type FeedbackNotification struct {
	BusinessId int `json:"businessId"`
	FeedbackId int `json:"feedbackId"`
}
