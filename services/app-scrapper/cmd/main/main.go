package main

import (
	"context"
	"github.com/oyamo/kplc-outage-microservice/services/app-scrapper/internal/core/scrapper/repositories"
	"github.com/oyamo/kplc-outage-microservice/services/app-scrapper/internal/core/scrapper/usecase"
	"github.com/oyamo/kplc-outage-microservice/services/app-scrapper/internal/worker"
	"github.com/oyamo/kplc-outage-microservice/services/app-scrapper/pkg/db"
	"github.com/oyamo/kplc-outage-microservice/services/app-scrapper/pkg/rmq"
	"github.com/qiniu/qmgo"
	"github.com/streadway/amqp"
	"log"
	"os"
	"os/signal"
	"time"
)

const (
	TaskInterval = time.Minute * 5
)

func main() {
	// Mongodb
	mongoHost := os.Getenv("MONGO_HOST")
	mongoPort := os.Getenv("MONGO_PORT")
	mongoDatabase := os.Getenv("MONGO_DATABASE")
	mongoDbUser := os.Getenv("MONGO_USER")
	mongoDbPassword := os.Getenv("MONGO_PASSWORD")

	// RabbitMQ
	rabbitHost := os.Getenv("RABBIT_HOST")
	rabbitUser := os.Getenv("RABBIT_USER")
	rabbitPassword := os.Getenv("RABBIT_PASS")
	rabbitPort := os.Getenv("RABBIT_PORT")

	var mongoClient *qmgo.Client
	var rmqConn *amqp.Connection
	var err error
	var config worker.Config

	if mongoHost == "" || mongoPort == "" || mongoDatabase == "" || mongoDbUser == "" || mongoDbPassword == "" {
		log.Fatalf("provide all env")
	}

	rmqConn, err = rmq.NewRmqClient(rabbitUser, rabbitPassword, rabbitHost, rabbitPort)
	if err != nil {
		log.Fatalf("Cannot connect to msg queue")
	} else {
		log.Println("Connected to messagequeue")
	}

	mongoClient, err = db.NewMgoClient(mongoHost, mongoPort, mongoDbUser, mongoDbPassword, mongoDatabase)
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
	rmqRepo := repositories.NewRmqRepo(rmqConn)

	scrapperUseCase := usecase.NewUseCase(mongodbRepo, webRepository, rmqRepo)

	// schedule the task
	worker.ScheduleWorker(TaskInterval, scrapperUseCase)

	var exit chan os.Signal = make(chan os.Signal, 1)
	signal.Notify(exit, os.Interrupt, os.Kill)
	<-exit
}
