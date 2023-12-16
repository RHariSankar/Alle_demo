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
	ChatGPTClient chatgpt.ChatGPTClient
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
	chatController := controllers.GetChatControllerInstance()
	chatController.NewChat(&chatRequest)
	response, err := ch.ChatGPTClient.ChatCompletion(chatRequest.Data)
	if err != nil {
		log.Printf("couldn't get chatgpt response %s", err)
		writer.WriteHeader(http.StatusInternalServerError)
	}
	responseTextChat := chat.TextChat{
		Type:     "chat",
		Role:     "system",
		Data:     response,
		DateTime: time.Now().Format(time.RFC3339Nano),
	}
	err = json.NewEncoder(writer).Encode(&responseTextChat)
	if err != nil {
		log.Printf("couldn't convert response to json %s", err)
		writer.WriteHeader(http.StatusInternalServerError)
	}
	chatController.NewChat(&responseTextChat)

}

func (ch *ChatHandler) AllChat(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(http.StatusOK)
	chatController := controllers.GetChatControllerInstance()
	chats, _ := chatController.AllChat()
	log.Println(len(chats))
	err := json.NewEncoder(writer).Encode(chats)
	if err != nil {
		log.Printf("couldn't marshal output %s", err)
		writer.WriteHeader(http.StatusInternalServerError)
	}
}
