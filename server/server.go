package server

import (
	"container/list"
	"context"
	"errors"
	"fmt"
	"log"
	"net"
	"sync"
	"time"

	grpc "google.golang.org/grpc"
)

// Errors.
var (
	ErrDuplicateQueue = errors.New("wolf: duplicate queue name")
	ErrNotFoundQueue  = errors.New("wolf: queue not found")
)

// Run register and run grpc server.
func Run(address string) (*grpc.Server, net.Listener, error) {
	fmt.Println("starting the grpc server")

	gs := grpc.NewServer()
	srv := &Server{}
	RegisterServerServer(gs, srv)

	l, err := net.Listen("tcp", address)
	if err != nil {
		return nil, nil, fmt.Errorf("listen %v", err)
	}

	go func() {
		err := gs.Serve(l)
		if err != nil {
			log.Fatal(err)
		}
	}()

	return gs, l, nil
}

// it's a queue pool.
var queues = make(map[string]Queue)

// Queue represents one queue struct.
type Queue struct {
	queueName string
	messages  *list.List
}

// Server represents message broker server.
type Server struct {
	wg sync.WaitGroup
}

// NewQueue try to create a new queue in the queues pool.
func (s *Server) NewQueue(ctx context.Context, in *NewQueueRequest) (*NewQueueResponse, error) {
	fmt.Printf("creating new queue with name: %s\n", in.QueueName)

	if _, ok := queues[in.QueueName]; ok {
		return nil, ErrDuplicateQueue
	}

	queue := Queue{queueName: in.QueueName, messages: list.New()}
	queues[in.QueueName] = queue

	fmt.Printf("created new queue with name: %s\n", in.QueueName)

	return &NewQueueResponse{QueueName: in.QueueName, Result: true}, nil
}

// PublishMessage try to enqueue a message in a queue.
func (s *Server) PublishMessage(ctx context.Context, in *PublishRequest) (*PublishResponse, error) {
	fmt.Printf("publish a message with id: %d to a queue: %s\n", in.Message.Id, in.QueueName)

	if _, ok := queues[in.QueueName]; !ok {
		return nil, ErrNotFoundQueue
	}

	queue := queues[in.QueueName]
	queue.messages.PushFront(in.Message)

	fmt.Printf("the message with id: %d enqueued to the queue: %s\n", in.Message.Id, in.QueueName)

	return &PublishResponse{QueueName: in.QueueName, Result: true}, nil
}

// Subscribe subscribe to the queue and getting the messages from her.
func (s *Server) Subscribe(in *SubscribeRequest, srv Server_SubscribeServer) error {
	fmt.Printf("subscribing to the queue %s\n", in.QueueName)

	if _, ok := queues[in.QueueName]; !ok {
		return ErrNotFoundQueue
	}

	queue := queues[in.QueueName]

	s.wg.Add(1)
	go s.Worker(&queue, srv)
	s.wg.Wait()

	return nil
}

// Worker check messages in queue and send to client.
func (s *Server) Worker(queue *Queue, srv Server_SubscribeServer) {
	// temp imitate unsibscribe
	var exit = make(chan (struct{}))
	time.AfterFunc(time.Second*10, func() { exit <- struct{}{} })

	for {
		select {
		case <-exit:
			fmt.Println("Exit")
			s.wg.Done()
			return
		default:
			if queue.messages.Len() > 0 {
				el := queue.messages.Back()
				val := el.Value

				message, ok := val.(*Message)
				if ok {
					resp := &SubscribeResponse{QueueName: queue.queueName, Message: message}

					if err := srv.Send(resp); err != nil {
						fmt.Println(err)
					}
				}

				queue.messages.Remove(el)
			}
		}
	}

}

func (s *Server) Unsubscribe() {}
