package cmd

import (
	"log"
	"net/http"
	"time"

	"github.com/coranlabs/CORAN_NWDAF/Application_entity/logger"
	analytics "github.com/coranlabs/CORAN_NWDAF/Application_entity/server/analytics"
	engine "github.com/coranlabs/CORAN_NWDAF/Application_entity/server/engine"
	events "github.com/coranlabs/CORAN_NWDAF/Application_entity/server/events"
	nbi_ml "github.com/coranlabs/CORAN_NWDAF/Application_entity/server/nbiml"
	sbi "github.com/coranlabs/CORAN_NWDAF/Application_entity/server/sbi"
	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
	"github.com/sirupsen/logrus"
)

type MainConfig struct {
	Server struct {
		// Addr string `envconfig:"SERVER_ADDR"`
	}
}

// ------------------------------------------------------------------------------
func Action() {
	logger.InitializeLogger(logrus.InfoLevel)
	// load the environment variables from the file .env
	err := godotenv.Load("config/nwdaf.env")
	if err != nil {
		log.Fatal("Error loading .env file sbi")
	}
	var config MainConfig
	err = envconfig.Process("", &config)
	if err != nil {
		log.Fatal(err.Error())
	}
	// Initialize internal package
	sbi.InitConfig()
	events.InitConfig()
	analytics.InitConfig()
	engine.InitConfig()

	// events
	IndividualNWDAFEventSubscriptionTransferDocumentApiService := events.NewIndividualNWDAFEventSubscriptionTransferDocumentApiService()
	IndividualNWDAFEventSubscriptionTransferDocumentApiController := events.NewIndividualNWDAFEventSubscriptionTransferDocumentApiController(
		IndividualNWDAFEventSubscriptionTransferDocumentApiService,
	)
	IndividualNWDAFEventsSubscriptionDocumentApiService := events.NewIndividualNWDAFEventsSubscriptionDocumentApiService()
	IndividualNWDAFEventsSubscriptionDocumentApiController := events.NewIndividualNWDAFEventsSubscriptionDocumentApiController(
		IndividualNWDAFEventsSubscriptionDocumentApiService,
	)
	NWDAFEventSubscriptionTransfersCollectionApiService := events.NewNWDAFEventSubscriptionTransfersCollectionApiService()
	NWDAFEventSubscriptionTransfersCollectionApiController := events.NewNWDAFEventSubscriptionTransfersCollectionApiController(
		NWDAFEventSubscriptionTransfersCollectionApiService,
	)
	NWDAFEventsSubscriptionsCollectionApiService := events.NewNWDAFEventsSubscriptionsCollectionApiService()
	NWDAFEventsSubscriptionsCollectionApiController := events.NewNWDAFEventsSubscriptionsCollectionApiController(
		NWDAFEventsSubscriptionsCollectionApiService,
	)

	// analytics
	NWDAFAnalyticsDocumentApiService := analytics.NewNWDAFAnalyticsDocumentApiService()
	NWDAFAnalyticsDocumentApiController := analytics.NewNWDAFAnalyticsDocumentApiController(
		NWDAFAnalyticsDocumentApiService,
	)
	NWDAFContextDocumentApiService := analytics.NewNWDAFContextDocumentApiService()
	NWDAFContextDocumentApiController := analytics.NewNWDAFContextDocumentApiController(
		NWDAFContextDocumentApiService,
	)

	//nbiml
	IndividualNWDAFMLModelProvisionSubscriptionDocumentApiService := nbi_ml.NewIndividualNWDAFMLModelProvisionSubscriptionDocumentApiService()
	IndividualNWDAFMLModelProvisionSubscriptionDocumentApiController := nbi_ml.NewIndividualNWDAFMLModelProvisionSubscriptionDocumentApiController(
		IndividualNWDAFMLModelProvisionSubscriptionDocumentApiService,
	)
	SubscriptionsCollectionApiService := nbi_ml.NewSubscriptionsCollectionApiService()
	SubscriptionsCollectionApiController := nbi_ml.NewSubscriptionsCollectionApiController(
		SubscriptionsCollectionApiService,
	)

	// Create routers for each section
	sbiRouter := sbi.NewRouter()
	eventsRouter := events.NewRouter(
		IndividualNWDAFEventSubscriptionTransferDocumentApiController,
		IndividualNWDAFEventsSubscriptionDocumentApiController,
		NWDAFEventSubscriptionTransfersCollectionApiController,
		NWDAFEventsSubscriptionsCollectionApiController,
	)
	analyticsRouter := analytics.NewRouter(
		NWDAFAnalyticsDocumentApiController,
		NWDAFContextDocumentApiController,
	)
	engineRouter := engine.NewRouter()
	nbimlRouter := nbi_ml.NewRouter(
		IndividualNWDAFMLModelProvisionSubscriptionDocumentApiController,
		SubscriptionsCollectionApiController,
	)

	// Create separate servers for each router
	go startServer("0.0.0.0:8080", sbiRouter)
	go startServer("0.0.0.0:8081", eventsRouter)
	go startServer("0.0.0.0:8082", analyticsRouter)
	go startServer("0.0.0.0:8083", engineRouter)
	go startServer("0.0.0.0:8084", nbimlRouter)

	// Block the main goroutine to keep the servers running
	select {}
}

// Helper function to start a server for a specific router
func startServer(addr string, handler http.Handler) {
	server := &http.Server{
		Addr:         addr,
		Handler:      handler,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	logger.EventLog.Infof("Server listening at %s", addr)
	log.Fatal(server.ListenAndServe())
}
