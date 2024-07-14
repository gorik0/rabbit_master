package main

import (
	"context"
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
	conn, err := (internal.NewRabbiConnection())
	HandlerErr(err, "connection error")
	client, err := internal.NewRabbitClient(conn)
	HandlerErr(err, "client error")
	defer client.Close()

	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, time.Second*15)
	defer cancel()
	msgBus, err := client.ConsumeMessages(ctx, "private_q", "commando", false)
	HandlerErr(err, "message bus error")

	g, ctx := errgroup.WithContext(ctx)

	cha := make(chan struct{})
	g.SetLimit(10)

	go func() {

		for mes := range msgBus {
			g.Go(func() error {

				log.Println("Received message: ", string(mes.Body))
				time.Sleep(time.Second * 10)
				err := mes.Ack(false)
				if err != nil {
					log.Println("Error acknowledging message: ", err)
					return err
				}
				log.Println("Acknowledged message: ")
				return nil
			})
		}

	}()
	<-cha
	//client.CreateQueue("private_q", true, false)
	//HandlerErr(client.QueueBind("private_q", "private.*", "army_event"), "queue1 binding error")
	//HandlerErr(client.QueueBind("private_q", "private.seaman.*", "army_event"), "queue1 binding error")
}
