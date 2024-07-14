package main

import (
	"context"
	rabbit "github.com/rabbitmq/amqp091-go"
	"log"
	"rabbit-master/internal"
	"strconv"
	"time"
)

//:::: RPC *2 connection both on consumer and producer*
//::: 			kind of exchange - direct

func main() {

	panamaChannel := make(chan struct{})

	//:: CONN for produce
	err, clientSend := setupSendClient()
	defer clientSend.Close()

	//:: CONN for recieve response

	clientRcv := setupRcvResponseClient(err)

	defer clientRcv.Close()

	//:::: WORK On RCV response

	//: queue & bind
	queue, err := clientRcv.CreateQueue("", true, false)
	HandlerErr(clientRcv.QueueBind(queue.Name, queue.Name, "army_callback"), "queue1 binding error")

	//: consume msg
	ctx2 := context.Background()
	msgBus, err := clientRcv.ConsumeMessages(ctx2, queue.Name, "consume_response", true)
	go func() {
		for d := range msgBus {
			log.Printf("Response to : %s", d.CorrelationId)
		}
	}()

	//:::: WORK On SEND
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	for i := 0; i < 10; i++ {

		err := clientSend.SendMessage(ctx, "army_event", "army_created", rabbit.Publishing{
			ContentType:   "text/plain",
			Body:          []byte("Hello World"),
			DeliveryMode:  rabbit.Persistent,
			ReplyTo:       queue.Name,
			CorrelationId: strconv.Itoa(i),
		})
		if err != nil {

			HandlerErr(err, "Error while send messs")
		}

	}

	<-panamaChannel
	time.Sleep(10 * time.Second)

}

func setupRcvResponseClient(err error) *internal.RabbitClient {
	connRcv, err := internal.NewRabbiConnection()
	HandlerErr(err, "error to create connection")

	clientRcv, err := internal.NewRabbitClient(connRcv)
	HandlerErr(err, "error to create client")
	return clientRcv
}

func setupSendClient() (error, *internal.RabbitClient) {
	connSend, err := internal.NewRabbiConnection()
	HandlerErr(err, "error to create connection")

	clientSend, err := internal.NewRabbitClient(connSend)
	HandlerErr(err, "error to create client")
	return err, clientSend
}

func HandlerErr(err error, message string) {
	if err != nil {
		log.Fatalf("Error: %s. Message: %s", err, message)
	}
}
