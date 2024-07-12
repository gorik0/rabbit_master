package main

import (
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

	err = client.CreateQueue("test1", true, false)
	if err != nil {
		panic(err)

	}
	err = client.CreateQueue("test2", false, true)
	if err != nil {
		panic(err)
	}
	time.Sleep(10 * time.Second)
}
