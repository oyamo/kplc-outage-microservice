package main

import (
	"context"
	"github.com/oyamo/kplc-outage-microservice/services/app-scrapper/internal/local/scrapper/repositories"
	"github.com/oyamo/kplc-outage-microservice/services/app-scrapper/internal/local/scrapper/usecase"
	"github.com/oyamo/kplc-outage-microservice/services/app-scrapper/internal/worker"
	"github.com/oyamo/kplc-outage-microservice/services/app-scrapper/pkg/databaseprovider"
	"github.com/qiniu/qmgo"
	"log"
	"os"
	"os/signal"
	"time"
)

const (
	TaskInterval = time.Minute * 5
)

func main() {
	mongoHost := os.Getenv("MONGO_HOST")
	mongoPort := os.Getenv("MONGO_PORT")
	mongoDatabase := os.Getenv("MONGO_DATABASE")
	mongoDbUser := os.Getenv("MONGO_USER")
	mongoDbPassword := os.Getenv("MONGO_PASSWORD")

	var mongoClient *qmgo.Client
	var err error
	var config worker.Config

	if mongoHost == "" || mongoPort == "" || mongoDatabase == "" || mongoDbUser == "" || mongoDbPassword == "" {
		log.Fatalf("provide all env")
	}

	mongoClient, err = databaseprovider.NewMgoClient(mongoHost, mongoPort, mongoDbUser, mongoDbPassword, mongoDatabase)
	if err != nil {
		log.Fatalf("cannot connect to %s: %s", mongoHost, err)
	} else {
		log.Printf("successfully connected to mongodb://%s:%s", mongoHost, mongoPort)
	}

	defer mongoClient.Close(context.Background())

	config.Database = mongoClient.Database(mongoDatabase)

	// Inititise repositories
	webRepository := repositories.NewWebRepository()
	mongodbRepo := repositories.NewMongoRepo(config.Database)
	scrapperUseCase := usecase.NewUseCase(mongodbRepo, webRepository)

	// schedule the task
	worker.ScheduleWorker(TaskInterval, scrapperUseCase)

	var exit chan os.Signal = make(chan os.Signal, 1)
	signal.Notify(exit, os.Interrupt, os.Kill)
	<-exit
}
