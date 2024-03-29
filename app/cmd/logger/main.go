package main

import (
	"fmt"
	"log"

	"github.com/allokate-ai/events/app/internal/config"
	"github.com/allokate-ai/events/app/internal/db"
	"github.com/allokate-ai/events/app/internal/queue"
	events "github.com/allokate-ai/events/app/pkg/client"
)

func main() {
	// Load the app's configuration settings.
	config, err := config.Get()
	if err != nil {
		log.Fatal(err)
	}

	// Connect to the database.
	pg, err := db.Connect(config.Database.Host, config.Database.Port, config.Database.User, config.Database.Password, config.Database.Database)
	if err != nil {
		log.Fatal(err)
	}
	defer pg.Close()

	// Initialize the database.
	if err := db.Init(pg); err != nil {
		log.Fatal(err)
	}

	amqp, err := queue.Connect(config.AMQPConfig.Host, config.AMQPConfig.Port, config.AMQPConfig.Username, config.AMQPConfig.Password)
	if err != nil {
		log.Fatal(err)
	}

	if err := queue.Init(amqp); err != nil {
		log.Fatal(err)
	}

	ch, err := amqp.Channel()
	if err != nil {
		log.Fatal(err)
	}

	// Run an endless consumer loop.
	if err := queue.Consume(ch, func(event *events.GenericEvent) bool {
		fmt.Println("Received event:", event)
		if _, err := db.UpsertEvent(pg, *event); err != nil {
			log.Println(err)
			return false
		}
		return true
	}); err != nil {
		log.Fatal(err)
	}

	log.Print("bye!")
}
