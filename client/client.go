package client

import (
	"fmt"
	"io"

	"github.com/Mycunycu/wolfmq-server/server"
	"google.golang.org/grpc"
)

// Client represents the client struct
type Client struct {
	server.ServerClient
}

// New connect then return a client instance.
func New(address string) (Client, error) {
	conn, err := connect(address)
	if err != nil {
		return Client{}, fmt.Errorf("client new connect %v", err)
	}

	cl := server.NewServerClient(conn)

	return Client{cl}, nil
}

func connect(address string) (*grpc.ClientConn, error) {
	return grpc.Dial(address, grpc.WithInsecure())
}

// ProcessMessage handle the received messages from subscribed queue.
func (c *Client) ProcessMessage(sub server.Server_SubscribeClient) {
	go func() {
		for {
			mess, err := sub.Recv()
			if err == io.EOF {
				fmt.Println("Unsubscribe")
				return
			}
			if err != nil {
				fmt.Printf("error subscription get message %v\n", err)
			}

			fmt.Println("[RECEIVED MESSAGE]:", mess.Message.Text)
		}
	}()
}
