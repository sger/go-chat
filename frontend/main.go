package main

import (
	"fmt"
	"net/http"

	"github.com/sger/go-chat/backend/server"
	"github.com/sger/go-chat/frontend/handlers"
)

func main() {
	go server.StartServer("0.0.0.0:50051")
	http.HandleFunc("/v1/chat", handlers.ChatHandler)
	http.Handle("/", http.StripPrefix("/", http.FileServer(http.Dir("./static"))))
	fmt.Println("Listening for connections on port", 8001)
	http.ListenAndServe(fmt.Sprintf(":%v", 8001), nil)
}
