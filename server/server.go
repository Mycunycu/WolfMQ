package server

import (
	"context"
	"fmt"
	"log"
	"net"

	"github.com/Mycunycu/wolfmq-server/config"
	grpc "google.golang.org/grpc"
)

func Run(cfg *config.Config) (*grpc.Server, net.Listener, error) {
	s := grpc.NewServer()
	srv := &Server{}
	RegisterServerServer(s, srv)

	address := fmt.Sprintf("%s:%s", cfg.Server.Host, cfg.Server.Port)
	l, err := net.Listen("tcp", address)
	if err != nil {
		return nil, nil, fmt.Errorf("listen %v", err)
	}

	go func() {
		err := s.Serve(l)
		if err != nil {
			log.Fatal(err)
		}
	}()

	return s, l, nil
}

type Server struct{}

func (s *Server) DoSomething(ctx context.Context, req *Request) (*Response, error) {
	return &Response{Message: req.Message + " " + "Hey!"}, nil
}
