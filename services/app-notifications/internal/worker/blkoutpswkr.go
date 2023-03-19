package worker

import (
	"context"
	"github.com/oyamo/kplc-outage-microservice/services/app-notifications/internal/core/subscription"
	"github.com/streadway/amqp"
)

const (
	ExchangeName = "newblkout"
	QueueName
	BlackoutIdKey = "blackout"
)

type blkoutpsworker struct {
	Amqpconn *amqp.Connection
	su       subscription.UseCase
}

func (w *blkoutpsworker) Run(errChan chan error, exitChan chan int, ctx context.Context) {
	// open achannel
	channel, err := w.Amqpconn.Channel()
	if err != nil {
		errChan <- err
		exitChan <- -1
		return
	}

	defer channel.Close()
	// declare a queue
	messageChan, err := channel.Consume(
		QueueName,
		"", true, // auto ack
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

		case <-ctx.Done():
			return
		}
	}
}

type BkoutpsWorker interface {
	Run(errChan chan error, exitChan chan int, ctx context.Context)
}

func NewWorker(amqp *amqp.Connection, su subscription.UseCase) BkoutpsWorker {
	return &blkoutpsworker{
		Amqpconn: amqp,
		su:       su,
	}
}
