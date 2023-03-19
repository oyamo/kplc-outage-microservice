package main

import (
	"context"
	"github.com/oyamo/kplc-outage-microservice/services/app-notifications/internal/core/subscription/repositories"
	"github.com/oyamo/kplc-outage-microservice/services/app-notifications/internal/core/subscription/usecase"
	"github.com/oyamo/kplc-outage-microservice/services/app-notifications/internal/server"
	"github.com/oyamo/kplc-outage-microservice/services/app-notifications/pkg/db"
	"github.com/qiniu/qmgo"
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
	grpcPort := os.Getenv("PORT")

	var mongoClient *qmgo.Client
	var err error
	var serverConfig server.Config

	if mongoHost == "" || mongoPort == "" || mongoDatabase == "" || mongoDbUser == "" || mongoDbPassword == "" ||
		grpcPort == "" {
		log.Fatalf("provide all env")
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
	usecaseI := usecase.NewUsecase(mgoRepo)

	// run the workers
	grpcExitChan := make(chan int, 1)
	errChan := make(chan error)
	quitChan := make(chan os.Signal)

	// RegisterKill signal
	signal.Notify(quitChan, os.Interrupt, syscall.SIGTERM)

	// Run GRPC Server
	server := server.NewServer(serverConfig, usecaseI)
	go server.Run(errChan, grpcExitChan)

	// Listen for new blackouts

run:
	for {
		select {
		case status, ok := <-grpcExitChan:
			{
				if !ok {
					break run
				}

				log.Printf("server exit with exit status: %d", status)
				close(errChan)
				close(grpcExitChan)
				close(quitChan)
			}
		case _, ok := <-quitChan:
			{
				if !ok {
					break run
				}
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
