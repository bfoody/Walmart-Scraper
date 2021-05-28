package main

import (
	"fmt"

	"github.com/bfoody/Walmart-Scraper/communication"
	"github.com/bfoody/Walmart-Scraper/logging"
	"go.uber.org/zap"
)

func main() {
	// Initialize logging.
	log, err := logging.Configure()
	if err != nil {
		fmt.Println("Error initializing logging: ", err)
	}

	log.Info("Hello world")

	conn, err := communication.ConnectAMQP("amqp://localhost:5672")
	if err != nil {
		log.Fatal(err.Error())
	}

	e := communication.NewQueueConnection(conn)

	q := "test"

	e.RegisterStatusUpdateHandler(func(su *communication.StatusUpdate) {
		log.Info("status update received", zap.String("statusUpdate", fmt.Sprintf("%+v", su)))
	})

	err = e.Consume(q)
	if err != nil {
		log.Fatal(err.Error())
	}

	forever := make(chan bool)
	<-forever
}
