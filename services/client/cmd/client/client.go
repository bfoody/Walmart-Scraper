package main

import (
	"fmt"
	"os"
	"os/signal"

	"syscall"

	"github.com/bfoody/Walmart-Scraper/communication"
	"github.com/bfoody/Walmart-Scraper/identity"
	"github.com/bfoody/Walmart-Scraper/logging"
	"github.com/bfoody/Walmart-Scraper/services/client/internal/receiver"
	"github.com/bfoody/Walmart-Scraper/utils/uuid"
)

func main() {
	// Initialize logging.
	log, err := logging.Configure()
	if err != nil {
		fmt.Println("Error initializing logging: ", err)
	}

	// TODO: move all this stuff into another function

	id := uuid.Generate()
	identity := identity.NewClient(id)

	conn, err := communication.ConnectAMQP("amqp://localhost:5672")
	if err != nil {
		log.Fatal(err.Error())
	}

	q := "test"
	e := communication.NewQueueConnection(conn, q)

	receiver := receiver.New(identity, log, e)

	err = e.Consume()
	if err != nil {
		log.Fatal(err.Error())
	}

	err = receiver.Start()
	if err != nil {
		log.Fatal(err.Error())
	}

	e.SendMessage(communication.StatusUpdate{
		FanoutPacket:     communication.FanoutPacket{SenderID: identity.ID},
		AvailableForWork: true,
	})

	log.Info(fmt.Sprintf("hello world! client initialized successfully as server %s", identity.ID))

	// Handle graceful shutdowns.
	s := make(chan os.Signal, 1)
	signal.Notify(s, os.Interrupt)
	signal.Notify(s, syscall.SIGTERM)
	go func() {
		<-s
		err = receiver.Shutdown()
		if err != nil {
			log.Fatal(err.Error())
		}
		os.Exit(0)
	}()

	forever := make(chan bool)
	<-forever
}
