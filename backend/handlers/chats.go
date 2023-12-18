package handlers

import (
	"alle/client/azure"
	chatgpt "alle/client/chatGPT"
	"alle/controllers"
	"alle/modals"
	"encoding/json"
	"log"
	"net/http"
	"time"
)

type ChatHandler struct {
	ChatGPTClient   chatgpt.ChatGPTClient
	ChatController  *controllers.ChatController
	AzureCLUClient  azure.AzureCLUClient
	ImageController controllers.ImageController
}

func (ch *ChatHandler) orchestrate(query string) ([]modals.Chat, error) {
	isIntent, entities, err := ch.AzureCLUClient.GetIntentAndEntity(query)
	log.Printf("azure clu returned %t %+v %s for query %s", isIntent, entities, err, query)
	if err != nil {
		return nil, err
	}
	if !isIntent {
		chatGptResponse, err := ch.ChatGPTClient.ChatCompletion(query)
		if err != nil {
			log.Printf("couldn't get chatgpt response %s", err)
			return nil, err
		}
		reply := modals.TextChat{
			Role:     "system",
			Text:     chatGptResponse,
			DateTime: time.Now().Format(time.RFC3339Nano),
		}
		reply.Type = reply.GetType()
		return []modals.Chat{&reply}, nil
	}
	return ch.ImageController.GetImagesByTags(entities)

}

func (ch *ChatHandler) handleNoEntitesError(err azure.AzureCLUNoEnitityError, writer http.ResponseWriter) {
	response := modals.TextChat{
		Role:     "system",
		DateTime: time.Now().Format(time.RFC3339),
		Text:     err.Error(),
	}
	response.Type = response.GetType()
	responses := []modals.Chat{&response}
	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(http.StatusOK)
	err1 := json.NewEncoder(writer).Encode(&responses)
	if err1 != nil {
		log.Printf("couldn't convert response to json %s", err)
		http.Error(writer, "couldn't convert response to json", http.StatusInternalServerError)
		return
	}
}

func (ch *ChatHandler) AddChat(writer http.ResponseWriter, request *http.Request) {

	var chatRequest modals.TextChat
	err := json.NewDecoder(request.Body).Decode(&chatRequest)
	if err != nil {
		log.Printf("couldn't convert request to object %s", err)
		writer.WriteHeader(http.StatusInternalServerError)
	}
	ch.ChatController.NewChat(&chatRequest)
	responses, err := ch.orchestrate(chatRequest.Text)
	if azureNoEntityError, ok := err.(*azure.AzureCLUNoEnitityError); ok {
		ch.handleNoEntitesError(*azureNoEntityError, writer)
		return
	}
	if err != nil {
		log.Printf("error in orchestrate  %s", err)
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}
	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(http.StatusOK)
	err = json.NewEncoder(writer).Encode(&responses)
	if err != nil {
		log.Printf("couldn't convert response to json %s", err)
		http.Error(writer, "couldn't convert response to json", http.StatusInternalServerError)
		return
	}
	for _, response := range responses {
		ch.ChatController.NewChat(response)
	}

}

func (ch *ChatHandler) AllChat(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(http.StatusOK)
	chats, _ := ch.ChatController.AllChat()
	err := json.NewEncoder(writer).Encode(chats)
	if err != nil {
		log.Printf("couldn't marshal output %s", err)
		writer.WriteHeader(http.StatusInternalServerError)
	}
}
