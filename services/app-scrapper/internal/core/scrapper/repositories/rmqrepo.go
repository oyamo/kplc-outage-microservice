package repositories

import (
	"github.com/oyamo/kplc-outage-microservice/services/app-scrapper/internal/core/scrapper"
	"github.com/streadway/amqp"
)

const (
	ExchangeName = "newblkout"
	QueueName
	BlackoutIdKey = "blackout"
)

type rmqrepo struct {
	conn *amqp.Connection
}

func (r rmqrepo) PublishId(id string) error {
	// open achannel
	channel, err := r.conn.Channel()
	if err != nil {
		return err
	}

	defer channel.Close()
	// declare a queue
	_, err = channel.QueueDeclare(
		QueueName,
		false, false, false, false, nil)
	if err != nil {
		return err
	}
	// publish a message
	err = channel.Publish(
		"",
		BlackoutIdKey,
		false,
		false,
		amqp.Publishing{ContentType: "text/plain", Body: []byte(id)},
	)

	if err != nil {
		return err
	}
	return nil
}

func NewRmqRepo(conn *amqp.Connection) scrapper.RmqRepo {
	return &rmqrepo{
		conn: conn,
	}
}
