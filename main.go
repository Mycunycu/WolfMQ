package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

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

	srv, lnr, err := server.Run(cfg)
	if err != nil {
		return fmt.Errorf("server run %v", err)
	}

	exit := make(chan os.Signal, 1)
	signal.Notify(exit, syscall.SIGTERM, syscall.SIGINT)
	<-exit

	srv.GracefulStop()
	lnr.Close()
	fmt.Println("Graceful shutdown")

	return nil
}
