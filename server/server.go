package server

import (
	"context"
)

func ConfigureOptions() (*Options, error) {
	opts := &Options{}

	return opts, nil
}

type Server struct{}

func (s *Server) DoSomething(ctx context.Context, req *Request) (*Response, error) {
	return &Response{Message: req.Message + " " + "Hey!"}, nil
}