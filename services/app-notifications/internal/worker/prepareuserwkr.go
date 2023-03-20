package worker

import (
	"context"
	"github.com/oyamo/kplc-outage-microservice/services/app-notifications/internal/core/subscription"
	"github.com/streadway/amqp"
)

type prepareuserwkr struct {
	suc      subscription.UseCase
	amqpConn *amqp.Connection
}

func (p prepareuserwkr) Run(errChan chan error, exitChan chan int, ctx context.Context) {
	// open achannel
	channel, err := p.amqpConn.Channel()
	if err != nil {
		errChan <- err
		exitChan <- -1
		return
	}

	defer channel.Close()
	// declare a queue
	messageChan, err := channel.Consume(
		"",
		QueueName,
		true,  // auto ack
		false, // exclusive
		false, // no local
		false, // no wait
		nil,   //args
	)
	if err != nil {
		errChan <- err
		exitChan <- -1
		return
	}

	for {
		select {
		case msg, ok := <-messageChan:
			if !ok {
				return
			}
			err := p.suc.ProvisionNextJobForUser(string(msg.Body))
			if err != nil {
				errChan <- err
			}
		case <-ctx.Done():
			return
		}
	}
}

type PrepareUserWorker interface {
	Run(errChan chan error, exitChan chan int, ctx context.Context)
}

func NewPrepareUserWorker(suc subscription.UseCase, amqpConn *amqp.Connection) PrepareUserWorker {
	return &prepareuserwkr{suc, amqpConn}
}
