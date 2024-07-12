package internal

import rabbit "github.com/rabbitmq/amqp091-go"

type RabbitClient struct {
	conn *rabbit.Connection
	ch   *rabbit.Channel
}

func NewRabbiConnection() (*rabbit.Connection, error) {
	return rabbit.Dial("amqp://gorik:1@localhost:5672/customers")
}

func NewRabbitClient(conn *rabbit.Connection) (*RabbitClient, error) {
	ch, err := conn.Channel()
	if err != nil {
		return nil, err
	}
	return &RabbitClient{conn, ch}, nil
}

func (rc *RabbitClient) Close() {
	rc.ch.Close()
}

func (rc *RabbitClient) CreateQueue(name string, durable, autoDelete bool) error {
	_, err := rc.ch.QueueDeclare(name, durable, autoDelete, false, false, nil)
	return err
}
