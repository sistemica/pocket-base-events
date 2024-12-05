package main

import (
	"log"
	"os"

	"github.com/pocketbase/pocketbase"
	"sistemica/pocket-engine/plugins"
)

func main() {
	app := pocketbase.New()

	redisUrl := os.Getenv("REDIS_URL")
	if redisUrl == "" {
		redisUrl = "localhost:6379"
	}

	redisPlugin := plugins.NewRedisPlugin(app, redisUrl)
	if redisPlugin == nil {
		log.Fatal("Failed to initialize Redis plugin")
	}

	redisPlugin.Register(app)
	redisPlugin.ListenAndProcessEvents(app)

	if err := app.Start(); err != nil {
		log.Fatal(err)
	}
}
