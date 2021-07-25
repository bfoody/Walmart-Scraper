package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/bfoody/Walmart-Scraper/communication"
	"github.com/bfoody/Walmart-Scraper/identity"
	"github.com/bfoody/Walmart-Scraper/logging"
	"github.com/bfoody/Walmart-Scraper/services/hub"
	"github.com/bfoody/Walmart-Scraper/services/hub/internal/database"
	"github.com/bfoody/Walmart-Scraper/services/hub/internal/database/postgres"
	"github.com/bfoody/Walmart-Scraper/services/hub/internal/service"
	"github.com/bfoody/Walmart-Scraper/services/hub/internal/supervisor"
	"github.com/bfoody/Walmart-Scraper/utils/uuid"
	"go.uber.org/zap"
)

func main() {
	// Initialize logging.
	log, err := logging.Configure()
	if err != nil {
		fmt.Println("Error initializing logging: ", err)
	}

	config, err := hub.LoadConfig()
	if err != nil {
		log.Fatal("unable to load config", zap.Error(err))
	}

	// TODO: move all this stuff into another function

	id := uuid.Generate()
	identity := identity.NewHub(id)

	db, err := postgres.Connect(postgres.ConnOptions{
		Host:           config.DatabaseURL,
		Port:           config.DatabasePort,
		DBName:         config.DatabaseName,
		Username:       config.DatabaseUsername,
		Password:       config.DatabasePassword,
		DisableSSLMode: true,
	})
	if err != nil {
		log.Fatal("error connecting to database", zap.Error(err))
	}

	productRepository := database.NewProductRepository(db)
	productInfoRepository := database.NewProductInfoRepository(db)
	productLocationRepository := database.NewProductLocationRepository(db)
	scrapeTaskRepository := database.NewScrapeTaskRepository(db)

	service := service.NewService(productRepository, productInfoRepository, productLocationRepository, scrapeTaskRepository)

	conn, err := communication.ConnectAMQP(config.AMQPURL)
	if err != nil {
		log.Fatal(err.Error())
	}

	e := communication.NewQueueConnection(conn, config.AMQPExchange)

	supervisor := supervisor.New(identity, log, e, service)

	err = e.Consume()
	if err != nil {
		log.Fatal(err.Error())
	}

	err = supervisor.Start()
	if err != nil {
		log.Fatal(err.Error())
	}

	log.Info(fmt.Sprintf("hello world! hub initialized successfully as hub %s", identity.ID))

	// Handle graceful shutdowns.
	s := make(chan os.Signal, 1)
	signal.Notify(s, os.Interrupt)
	signal.Notify(s, syscall.SIGTERM)
	go func() {
		<-s
		err = supervisor.Shutdown()
		if err != nil {
			log.Fatal(err.Error())
		}
		os.Exit(0)
	}()

	// Run forever (until stopped)
	forever := make(chan bool)
	<-forever
}
