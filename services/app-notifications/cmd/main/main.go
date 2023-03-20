package main

import (
	"context"
	"github.com/oyamo/kplc-outage-microservice/services/app-notifications/internal/core/subscription/repositories"
	"github.com/oyamo/kplc-outage-microservice/services/app-notifications/internal/core/subscription/usecase"
	"github.com/oyamo/kplc-outage-microservice/services/app-notifications/internal/server"
	"github.com/oyamo/kplc-outage-microservice/services/app-notifications/internal/worker"
	"github.com/oyamo/kplc-outage-microservice/services/app-notifications/pkg/db"
	"github.com/oyamo/kplc-outage-microservice/services/app-notifications/pkg/rmq"
	"github.com/qiniu/qmgo"
	"github.com/streadway/amqp"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
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

	// Grpc
	grpcPort := os.Getenv("PORT")

	var mongoClient *qmgo.Client
	var rmqConn *amqp.Connection
	var err error
	var serverConfig server.Config

	if mongoHost == "" || mongoPort == "" || mongoDatabase == "" || mongoDbUser == "" || mongoDbPassword == "" ||
		grpcPort == "" {
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

	serverConfig.Database = mongoClient.Database(mongoDatabase)

	// create repositories
	mgoRepo := repositories.NewMongoRepo(serverConfig.Database)
	rmqRepo := repositories.NewAmqpRepo(rmqConn)
	usecaseI := usecase.NewUsecase(mgoRepo, rmqRepo)

	// run the workers
	workerExitChan := make(chan int, 1)
	errChan := make(chan error)
	quitChan := make(chan os.Signal)
	cxt, cancelWork := context.WithCancel(context.Background())

	// RegisterKill signal
	signal.Notify(quitChan, os.Interrupt, syscall.SIGTERM, syscall.SIGQUIT)

	// Run GRPC Server
	server := server.NewServer(serverConfig, usecaseI)
	go server.Run(errChan, workerExitChan)

	// Listen for new blackouts
	blkoutWorker := worker.NewBlkoutWorker(rmqConn, usecaseI)
	go blkoutWorker.Run(errChan, workerExitChan, cxt)

	// Listen for new user additions
	prepUserWorker := worker.NewPrepareUserWorker(usecaseI, rmqConn)
	go prepUserWorker.Run(errChan, workerExitChan, cxt)

	// Notifications worker

run:
	for {
		select {
		case status, ok := <-workerExitChan:
			{
				if !ok {
					break run
				}

				log.Printf("server exit with exit status: %d", status)
				close(errChan)
				close(workerExitChan)
				close(quitChan)
			}
		case _, ok := <-quitChan:
			{
				if !ok {
					break run
				}
				cancelWork()
				log.Printf("program quiting")
				break run
			}

		case err, ok := <-errChan:
			{
				if !ok {
					break run
				}

				log.Printf("error: %s", err)
			}
		}
	}
}
