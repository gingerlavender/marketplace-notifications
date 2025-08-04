package main

import (
	"marketplace-notifications/internal/app"
)

func main() {
	app := app.NewApp()

	app.Run()
}
