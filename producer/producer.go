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
	for i := 0; i < 10; i++ {

		err := client.SendMessage(ctx, "army_event", "private.vasya", rabbit.Publishing{
			ContentType:  "text/plain",
			Body:         []byte("Hello World"),
			DeliveryMode: rabbit.Transient,
		})
		if err != nil {
			HandlerErr(err, "Error while send messs")
		}
		err = client.SendMessage(ctx, "army_event", "private.pyatro", rabbit.Publishing{
			ContentType:  "text/plain",
			Body:         []byte("Hello World"),
			DeliveryMode: rabbit.Persistent,
		})
		if err != nil {
			HandlerErr(err, "Error while send messs")
		}

	}

	time.Sleep(10 * time.Second)

}

func HandlerErr(err error, message string) {
	if err != nil {
		log.Fatalf("Error: %s. Message: %s", err, message)
	}
}
