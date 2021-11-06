package main

import (
	"fmt"
	"log"
	"net"

	"github.com/Mycunycu/wolfmq-server/server"
	"google.golang.org/grpc"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	s := grpc.NewServer()
	srv := &server.Server{}
	server.RegisterServerServer(s, srv)

	l, err := net.Listen("tcp", ":8825")
	if err != nil {
		return fmt.Errorf("listen %v", err)
	}

	if err := s.Serve(l); err != nil {
		return fmt.Errorf("serve %v", err)
	}

	return nil
}