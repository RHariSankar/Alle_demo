package handlers

import (
	"alle/chat"
	chatgpt "alle/client/chatGPT"
	"alle/controllers"
	"encoding/json"
	"log"
	"net/http"
	"time"
)

type ChatHandler struct {
	ChatGPTClient  chatgpt.ChatGPTClient
	ChatController *controllers.ChatController
}

func (ch *ChatHandler) AddChat(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(http.StatusOK)
	var chatRequest chat.TextChat
	err := json.NewDecoder(request.Body).Decode(&chatRequest)
	if err != nil {
		log.Printf("couldn't convert request to object %s", err)
		writer.WriteHeader(http.StatusInternalServerError)
	}
	ch.ChatController.NewChat(&chatRequest)
	// response, err := ch.ChatGPTClient.ChatCompletion(chatRequest.Text)
	response := "This is a sample response"
	if err != nil {
		log.Printf("couldn't get chatgpt response %s", err)
		writer.WriteHeader(http.StatusInternalServerError)
	}
	responseTextChat := chat.TextChat{
		Type:     "chat",
		Role:     "system",
		Text:     response,
		DateTime: time.Now().Format(time.RFC3339Nano),
	}
	responses := make([]chat.Chat, 0)
	responses = append(responses, &responseTextChat)
	err = json.NewEncoder(writer).Encode(&responses)
	if err != nil {
		log.Printf("couldn't convert response to json %s", err)
		writer.WriteHeader(http.StatusInternalServerError)
	}
	ch.ChatController.NewChat(&responseTextChat)

}

func (ch *ChatHandler) AllChat(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(http.StatusOK)
	chats, _ := ch.ChatController.AllChat()
	log.Println(len(chats))
	err := json.NewEncoder(writer).Encode(chats)
	if err != nil {
		log.Printf("couldn't marshal output %s", err)
		writer.WriteHeader(http.StatusInternalServerError)
	}
}
