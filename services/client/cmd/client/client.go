package main

import (
	"fmt"

	"github.com/bfoody/Walmart-Scraper/communication"
	"github.com/bfoody/Walmart-Scraper/logging"
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

	err = e.Consume(q)
	if err != nil {
		log.Fatal(err.Error())
	}

	e.SendMessage(q, communication.StatusUpdate{
		FanoutPacket:     communication.FanoutPacket{SenderID: "server1"},
		AvailableForWork: true,
	})

	forever := make(chan bool)
	<-forever
}
