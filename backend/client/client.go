package client

import (
	"context"
	"io"
	"log"

	"google.golang.org/grpc/grpclog"

	chat "github.com/sger/go-chat/backend/protos"
	"google.golang.org/grpc"
)

// Callback will be called on client
type Callback interface {
	NewMessage(message string)
}

// ChatClient creates a new client that can be used for sending and receriving messages
type ChatClient struct {
	client     chat.ChatClient
	callback   Callback
	connection *grpc.ClientConn
}

func Create(server string) *ChatClient {
	connection, err := grpc.Dial(server, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("server did not connect: %v", err)
	}

	client := chat.NewChatClient(connection)

	return &ChatClient{client: client, connection: connection}
}

func (c *ChatClient) WriteMessage(user string, message string) {
	c.client.Send(context.Background(), &chat.Message{From: user, Content: message})
}

func (c *ChatClient) ReadStream(callback Callback) {
	c.callback = callback
	go c.startRead()
}

func (c *ChatClient) startRead() {
	stream, err := c.client.Stream(context.Background(), &chat.Filter{})
	if err != nil {
		grpclog.Fatalf("%v.ListFeatures(_) = _, %v", c.client, err)
	}
	for {
		message, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			grpclog.Fatalf("%v.ListFeatures(_) = _, %v", c.client, err)
		}

		c.callback.NewMessage(message.Content)
	}
}

func (c *ChatClient) Close() {
	c.connection.Close()
}
