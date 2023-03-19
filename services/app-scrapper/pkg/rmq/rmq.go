package rmq

import (
	"fmt"
	"github.com/streadway/amqp"
)

func NewRmqClient(user, pass, host, port string) (*amqp.Connection, error) {
	uri := fmt.Sprintf("amqp://%s:%s@%s:%s",user, pass, host, port)
	conn, err := amqp.Dial(uri)

	if err != nil {
		return nil, err
	}
	return conn, nil
}
