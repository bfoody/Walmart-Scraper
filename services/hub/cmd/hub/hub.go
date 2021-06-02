package main

import (
	"fmt"

	"github.com/bfoody/Walmart-Scraper/communication"
	"github.com/bfoody/Walmart-Scraper/identity"
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

	// TODO: move all this stuff into another function

	id := uuid.Generate()
	identity := identity.NewHub(id)

	conn, err := communication.ConnectAMQP("amqp://localhost:5672")
	if err != nil {
		log.Fatal(err.Error())
	}

	q := "test"
	e := communication.NewQueueConnection(conn, q)

	supervisor := supervisor.New(identity, log, e)

	err = e.Consume()
	if err != nil {
		log.Fatal(err.Error())
	}

	err = supervisor.Start()
	if err != nil {
		log.Fatal(err.Error())
	}

	log.Info(fmt.Sprintf("hello world! hub initialized successfully as hub %s", identity.ID))

	// Run forever (until stopped)
	forever := make(chan bool)
	<-forever

	err = supervisor.Shutdown()
	if err != nil {
		log.Fatal(err.Error())
	}
}
