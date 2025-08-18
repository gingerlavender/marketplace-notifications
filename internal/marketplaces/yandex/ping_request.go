package yandex

import "time"

type PingRequest struct {
	NotificatonType string    `json:"notificationType"`
	Time            time.Time `json:"time"`
}
