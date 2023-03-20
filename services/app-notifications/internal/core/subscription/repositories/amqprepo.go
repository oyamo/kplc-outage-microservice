package repositories

import (
	"github.com/oyamo/kplc-outage-microservice/services/app-notifications/internal/core/subscription"
	"github.com/streadway/amqp"
)

const QueueName = "new-user"

type amqprepo struct {
	conn *amqp.Connection
}

func (a amqprepo) PublishUserId(userId string) error {

	// open achannel
	channel, err := a.conn.Channel()
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
		QueueName,
		false,
		false,
		amqp.Publishing{ContentType: "text/plain", Body: []byte(userId)},
	)

	if err != nil {
		return err
	}
	return nil
}

func NewAmqpRepo(conn *amqp.Connection) subscription.AmqpRepo {
	return &amqprepo{conn: conn}
}
