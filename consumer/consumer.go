package main

import (
	"context"
	"github.com/rabbitmq/amqp091-go"
	"golang.org/x/sync/errgroup"
	"log"
	"rabbit-master/internal"
	"time"
)

func HandlerErr(err error, message string) {
	if err != nil {
		log.Fatalf("Error: %s. Message: %s", err, message)
	}
}

func main() {

	panamaChannel := make(chan struct{})
	//:: CONN for consume
	connConsume, err := internal.NewRabbiConnection()
	HandlerErr(err, "error to create connection")

	clientConsume, err := internal.NewRabbitClient(connConsume)
	HandlerErr(err, "error to create client")

	defer clientConsume.Close()

	queue, err := clientConsume.CreateQueue("", true, false)
	HandlerErr(clientConsume.QueueBind(queue.Name, "", "army_event"), "queue1 binding error")

	//:: CONN for recieve response

	connSendResp, err := internal.NewRabbiConnection()
	HandlerErr(err, "error to create connection")

	clientSendResp, err := internal.NewRabbitClient(connSendResp)
	HandlerErr(err, "error to create client")

	defer clientSendResp.Close()

	//::::: WORK on consuming and sennding response

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	ctxResp := context.Background()
	msgBus, err := clientConsume.ConsumeMessages(ctx, queue.Name, "consumer", false)
	HandlerErr(err, "error to create consumer")
	g, ctx := errgroup.WithContext(ctx)
	go func() {

		println("start to consume. ... ")
		for mesg := range msgBus {
			println("start another go . ....")
			println(len(msgBus))
			msg := mesg
			g.Go(func() error {

				log.Println("GOtta --->> ", string(msg.Body))

				time.Sleep(time.Second * 1)
				err2 := msg.Ack(false)
				//log.Println(err2)
				HandlerErr(err2, "error to ack message")

				//	send respond to

				err2 = clientSendResp.SendMessage(ctxResp, "army_callback", msg.ReplyTo, amqp091.Publishing{
					ContentType:   "text/plain",
					DeliveryMode:  amqp091.Transient,
					CorrelationId: msg.CorrelationId,
				})
				if err2 != nil {
					println("error akcowledged mes")

				}
				println("end of go")
				return nil

			},
			)

		}
	}()
	<-panamaChannel
}
