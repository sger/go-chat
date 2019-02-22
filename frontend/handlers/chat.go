package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/sger/go-chat/backend/client"
)

// ChatRequest send chat request
type ChatRequest struct {
	Message string `json: "message"`
}

// ChatHandler handle Chat requests
func ChatHandler(rw http.ResponseWriter, r *http.Request) {
	var request ChatRequest
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&request)

	if err != nil || request.Message == "" || len(request.Message) > 154 {
		http.Error(rw, "Invalid message length", http.StatusBadRequest)
		return
	}

	c := client.Create("0.0.0.0:50051")
	defer c.Close()

	c.WriteMessage("user", request.Message)
}
