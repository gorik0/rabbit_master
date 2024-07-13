package main

import (
	"context"
	rabbit "github.com/rabbitmq/amqp091-go"
	"log"
	"rabbit-master/internal"
	"time"
)

func main() {
	conn, err := internal.NewRabbiConnection()
	if err != nil {
		panic(err)
	}
	client, err := internal.NewRabbitClient(conn)
	if err != nil {
		panic(err)
	}
	defer client.Close()
	//client.CreateQueue("private_q", true, false)
	//HandlerErr(client.QueueBind("private_q", "private.*", "army_event"), "queue1 binding error")
	//HandlerErr(client.QueueBind("private_q", "private.seaman.*", "army_event"), "queue1 binding error")

	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	HandlerErr(client.SendMessage(ctx, "army_event", "private.to_forest", rabbit.Publishing{
		ContentType:  "text/plain",
		DeliveryMode: rabbit.Persistent,
		Body:         []byte("Hello Forest!"),
	}), "message sending error")
	time.Sleep(10 * time.Second)

}

func HandlerErr(err error, message string) {
	if err != nil {
		log.Fatalf("Error: %s. Message: %s", err, message)
	}
}
