package internal

import (
	"context"
	rabbit "github.com/rabbitmq/amqp091-go"
	"log"
)

type RabbitClient struct {
	conn *rabbit.Connection
	ch   *rabbit.Channel
}

func NewRabbiConnection() (*rabbit.Connection, error) {
	return rabbit.Dial("amqp://gorik:1@localhost:5672/army")
}

func NewRabbitClient(conn *rabbit.Connection) (*RabbitClient, error) {
	ch, err := conn.Channel()
	if err != nil {
		return nil, err
	}

	err = ch.Confirm(false)
	if err != nil {
		panic(err)
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
func (rc *RabbitClient) QueueBind(name string, binding string, exchange string) error {
	return rc.ch.QueueBind(name, binding, exchange, false, nil)

}

func (rc *RabbitClient) SendMessage(ctx context.Context, exchange, key string, opts rabbit.Publishing) error {
	confirmed, err := rc.ch.PublishWithDeferredConfirmWithContext(ctx, exchange, key, true, false, opts)
	if err != nil {
		println("Msg hasn't been comfrmed!!!")
		return err
	}
	log.Println(confirmed.Wait())
	return nil

}

func (rc *RabbitClient) ConsumeMessages(ctx context.Context, queue, consumer string, autoAck bool) (<-chan rabbit.Delivery, error) {
	return rc.ch.ConsumeWithContext(ctx, queue, consumer, autoAck, false, false, false, nil)
}
