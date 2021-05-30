package main

import (
	"fmt"

	"github.com/bfoody/Walmart-Scraper/communication"
	"github.com/bfoody/Walmart-Scraper/logging"
	"github.com/bfoody/Walmart-Scraper/services/hub/internal/supervisor"
	"github.com/bfoody/Walmart-Scraper/utils/uuid"
)

func main() {
	// Initialize logging.
	log, err := logging.Configure()
	if err != nil {
		fmt.Println("Error initializing logging: ", err)
	}

	id := uuid.Generate()

	conn, err := communication.ConnectAMQP("amqp://localhost:5672")
	if err != nil {
		log.Fatal(err.Error())
	}

	e := communication.NewQueueConnection(conn)

	q := "test"

	supervisor := supervisor.New(log)

	e.RegisterStatusUpdateHandler(func(su *communication.StatusUpdate) {
		supervisor.PipeStatusUpdate(*su)
	})

	err = e.Consume(q)
	if err != nil {
		log.Fatal(err.Error())
	}

	err = supervisor.Start()
	if err != nil {
		log.Fatal(err.Error())
	}

	log.Info(fmt.Sprintf("hello world! hub initialized successfully as hub %s", id))

	forever := make(chan bool)
	<-forever

	err = supervisor.Shutdown()
	if err != nil {
		log.Fatal(err.Error())
	}
}
