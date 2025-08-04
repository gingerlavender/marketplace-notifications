package app

import (
	"fmt"
	"log"
	"marketplace-notifications/internal/client"
	"marketplace-notifications/internal/config"
	"marketplace-notifications/internal/monitor"
	"marketplace-notifications/internal/telegram"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

type App struct {
	config  *config.ServerConfig
	monitor *monitor.Monitor
}

func NewApp() *App {
	if err := godotenv.Load(); err != nil {
		log.Printf("Warning: Could not load .env file: %v", err)
		log.Println("Assuming environment variables are set directly")
	}

	config, err := config.Load()
	if err != nil {
		log.Fatal(err)
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
	router.GET("/start", app.start)
	router.GET("/stop", app.stop)

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
