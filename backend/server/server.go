package server

import (
	"context"
	"fmt"
	"log"
	"net"

	chat "github.com/sger/go-chat/backend/protos"
	"google.golang.org/grpc"
)

const port = ":50051"

type broadcaster struct {
	channels map[chat.Chat_StreamServer]chan chat.Message
}

func (b *broadcaster) Write(message chat.Message) {
	for _, channel := range b.channels {
		go func() {
			channel <- message
		}()
	}
}

func (b *broadcaster) Listen(stream chat.Chat_StreamServer) chan chat.Message {
	fmt.Println("New channel added")

	channel := make(chan chat.Message)
	b.channels[stream] = channel

	return channel
}

func newBroadcaster() *broadcaster {
	m := make(map[chat.Chat_StreamServer]chan chat.Message)
	return &broadcaster{channels: m}
}

type server struct {
	broadcaster *broadcaster
}

func (s *server) Send(ctx context.Context, in *chat.Message) (*chat.Error, error) {
	s.broadcaster.Write(*in)
	return &chat.Error{Code: -1}, nil
}

func (s *server) Stream(filter *chat.Filter, stream chat.Chat_StreamServer) error {
	channel := s.broadcaster.Listen(stream)

	for {
		message := <-channel
		if err := stream.Send(&message); err != nil {
			fmt.Println("Cannot send message")
			delete(s.broadcaster.channels, stream)
			return err
		}
	}
}

func StartServer(location string) {
	lis, err := net.Listen("tcp", location)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	rpcServer := grpc.NewServer()
	chatServer := &server{broadcaster: newBroadcaster()}

	chat.RegisterChatServer(rpcServer, chatServer)

	fmt.Println("Starting server on: ", port)

	rpcServer.Serve(lis)
}
