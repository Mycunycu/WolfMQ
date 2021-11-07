package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/Mycunycu/wolfmq-server/client"
	"github.com/Mycunycu/wolfmq-server/config"
	"github.com/Mycunycu/wolfmq-server/server"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	err := config.Init()
	if err != nil {
		return fmt.Errorf("config init %v", err)
	}

	cfg := config.Get()
	address := fmt.Sprintf("%s:%s", cfg.Server.Host, cfg.Server.Port)

	srv, lnr, err := server.Run(address)
	if err != nil {
		return fmt.Errorf("server run %v", err)
	}

	cl, err := client.New(address)
	if err != nil {
		return fmt.Errorf("run client new %v", err)
	}

	_, err = cl.NewQueue(context.Background(), &server.NewQueueRequest{QueueName: "ololo"})
	if err != nil {
		return fmt.Errorf("new queue %v", err)
	}

	msg := &server.Message{Id: 1, Text: "this is a message text"}
	_, err = cl.PublishMessage(context.Background(), &server.PublishRequest{QueueName: "ololo", Message: msg})
	if err != nil {
		return fmt.Errorf("publish message %v", err)
	}

	msg = &server.Message{Id: 2, Text: "one more message"}
	_, err = cl.PublishMessage(context.Background(), &server.PublishRequest{QueueName: "ololo", Message: msg})
	if err != nil {
		return fmt.Errorf("publish message %v", err)
	}

	subscr, err := cl.Subscribe(context.Background(), &server.SubscribeRequest{QueueName: "ololo"})
	if err != nil {
		return fmt.Errorf("subscribe error %v", err)
	}

	go cl.ProcessMessage(subscr)

	exit := make(chan os.Signal, 1)
	signal.Notify(exit, syscall.SIGTERM, syscall.SIGINT)
	<-exit

	fmt.Println("Graceful shutdown")
	srv.GracefulStop()
	lnr.Close()

	return nil
}
