package monitor

import (
	"context"
	"log"
	"marketplace-notifications/internal/client"
	"marketplace-notifications/internal/config"
	"marketplace-notifications/internal/telegram"
	"sync"
	"time"
)

type Monitor struct {
	mutex                sync.RWMutex
	isRunning            bool
	lastCheck            time.Time
	lastUpdateDiscovered time.Time
	config               *config.MonitorConfig
	apiClient            *client.APIClient
	notifier             *telegram.TelegramNotifier
	ctx                  context.Context
	cancel               context.CancelFunc
}

func NewMonitor(config *config.MonitorConfig, apiClent *client.APIClient, notifier *telegram.TelegramNotifier) *Monitor {
	return &Monitor{
		config:    config,
		apiClient: apiClent,
		notifier:  notifier,
	}
}

func (monitor *Monitor) Start() {
	monitor.mutex.Lock()
	defer monitor.mutex.Unlock()

	if monitor.isRunning {
		log.Println("Already running")
		return
	}

	monitor.ctx, monitor.cancel = context.WithCancel(context.Background())
	monitor.isRunning = true

	log.Println("Starting monitor...")
	log.Printf("Check interval: %s", monitor.config.CheckInterval)

	go monitor.run()
}

func (monitor *Monitor) Stop() {
	monitor.mutex.Lock()
	defer monitor.mutex.Unlock()

	if !monitor.isRunning {
		log.Println("Monitor is not running")
		return
	}

	monitor.isRunning = false

	if monitor.cancel != nil {
		monitor.cancel()
	}
}

func (monitor *Monitor) IsRunning() bool {
	monitor.mutex.Lock()
	defer monitor.mutex.Unlock()

	return monitor.isRunning
}

func (monitor *Monitor) GetInfo() map[string]any {
	monitor.mutex.RLock()
	defer monitor.mutex.RUnlock()

	return map[string]any{
		"isRunning":            monitor.isRunning,
		"lastCheck":            monitor.lastCheck,
		"lastUpdateDiscovered": monitor.lastUpdateDiscovered,
	}
}

func (monitor *Monitor) run() {
	ticker := time.NewTicker(monitor.config.CheckInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			monitor.checkForUpdates()
		case <-monitor.ctx.Done():
			log.Println("Monitor stopped")
			return
		}
	}
}

func (monitor *Monitor) checkForUpdates() {
	log.Println("Checking for questions...")

	questions, questionsErr := monitor.apiClient.FetchWBQuestionsForPeriod(monitor.config.CheckInterval)
	if questionsErr != nil {
		log.Printf("Failed to check for questions: %v", questionsErr)

	}

	log.Println("Checking for feedbacks...")

	feedbacks, feedbacksErr := monitor.apiClient.FetchWBFeedbacksForPeriod(monitor.config.CheckInterval)
	if feedbacksErr != nil {
		log.Printf("Failed to check for feedbacks: %v", feedbacksErr)
	}

	if questionsErr == nil && feedbacksErr == nil {
		monitor.lastCheck = time.Now()
	}

	newQuestionsNumber, newFeedbacksNumber := len(questions), len(feedbacks)

	if newQuestionsNumber > 0 || newFeedbacksNumber > 0 {
		monitor.lastUpdateDiscovered = monitor.lastCheck
	}

	log.Printf("Found %d new questions and %d new feedbacks", newQuestionsNumber, newFeedbacksNumber)

	for _, question := range questions {
		if err := monitor.notifier.SendQuestionNotification(question); err != nil {
			log.Printf("Failed to send notification for question with id %s: %v", question.Id, err)
		} else {
			log.Printf("Sent question notification with id %s", question.Id)
		}
	}

	for _, feedback := range feedbacks {
		if err := monitor.notifier.SendFeedbackNotification(feedback); err != nil {
			log.Printf("Failed to send notification for feedback with id %s: %v", feedback.Id, err)
		} else {
			log.Printf("Sent feedback notification with id %s", feedback.Id)
		}
	}

}
