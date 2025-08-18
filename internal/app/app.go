package app

import (
	"encoding/json"
	"fmt"
	"log"
	"marketplace-notifications/internal/client"
	"marketplace-notifications/internal/config"
	"marketplace-notifications/internal/marketplaces/yandex"
	"marketplace-notifications/internal/monitor"
	"marketplace-notifications/internal/telegram"
	"marketplace-notifications/internal/utils/ip"
	"net"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

type App struct {
	config  *config.ServerConfig
	monitor *monitor.Monitor
}

func NewApp() *App {
	log.SetOutput(os.Stdout)

	if err := godotenv.Load(); err != nil {
		log.Printf("[WARN] Warning: Could not load .env file: %v", err)
		log.Println("[INFO] Assuming environment variables are set directly")
	}

	config, err := config.Load()
	if err != nil {
		log.Fatal("[ERROR] Failed to load config: ", err)
	}

	apiClient := client.NewAPIClient(&config.API)
	notifier := telegram.NewTelegramNotifier(&config.Telegram)
	monitor := monitor.NewMonitor(&config.Monitor, apiClient, notifier)

	return &App{
		config:  &config.Server,
		monitor: monitor,
	}
}

func (app *App) Run() {
	router := gin.Default()

	router.GET("/info", app.getInfo)
	router.POST("/start", app.start)
	router.POST("/stop", app.stop)
	router.POST("/notify", app.handleYandexPing)

	router.Run(fmt.Sprintf(":%d", app.config.Port))
}

func (app *App) getInfo(c *gin.Context) {
	c.JSON(http.StatusOK, app.monitor.GetInfo())
}

func (app *App) start(c *gin.Context) {
	if token := c.Query("token"); token != app.config.ControlToken {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "missing on incorrect control token"})
		return
	}

	if app.monitor.IsRunning() {
		c.JSON(http.StatusOK, gin.H{"message": "already running"})
		return
	}

	app.monitor.Start()

	c.JSON(http.StatusOK, gin.H{"message": "started running"})
}

func (app *App) stop(c *gin.Context) {
	if token := c.Query("token"); token != app.config.ControlToken {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "missing on incorrect control token"})
		return
	}

	if !app.monitor.IsRunning() {
		c.JSON(http.StatusOK, gin.H{"message": "monitor is not running"})
		return
	}

	app.monitor.Stop()

	c.JSON(http.StatusOK, gin.H{"message": "running stops..."})
}

func (app *App) handleYandexPing(c *gin.Context) {
	clientIP := net.ParseIP(c.ClientIP())
	if clientIP == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "unable to determine client IP"})
		return
	}

	if !ip.IsAllowed(clientIP, yandex.IPWhitelist) {
		c.JSON(http.StatusForbidden, gin.H{"error": "IP is not in the whitelist"})
		return
	}

	var body yandex.PingRequest

	err := json.NewDecoder(c.Request.Body).Decode(&body)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "unable to decode request body"})
		return
	}

	if body.NotificatonType != "PING" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "wrong notification type"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"version": "0.1.0", "name": "Bezlimit", "time": body.Time.Format(time.RFC3339)})
}
