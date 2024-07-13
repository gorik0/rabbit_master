package main

import (
	"context"
	"fmt"
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
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	msgBus, err := client.ConsumeMessages(ctx, "private_q", "private.to_sea", true)
	HandlerErr(err, "consumer error")

	go func() {

		for msg := range msgBus {
			fmt.Println("Gotta msg!!!  --->>>  " + string(msg.Body))
		}
	}()
	time.Sleep(10 * time.Second)

	//client.CreateQueue("private_q", true, false)
	//HandlerErr(client.QueueBind("private_q", "private.*", "army_event"), "queue1 binding error")
	//HandlerErr(client.QueueBind("private_q", "private.seaman.*", "army_event"), "queue1 binding error")
}
