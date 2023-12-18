package chatgpt

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
)

type GPTMesssage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type ChatGPTRequest struct {
	Model   string        `json:"model"`
	Message []GPTMesssage `json:"messages"`
}

type Choice struct {
	Index        int         `json:"index"`
	Message      GPTMesssage `json:"message"`
	FinishReason string      `json:"finish_reason"`
}

type ChatGPTResponse struct {
	Choices []Choice `json:"choices"`
}

type ChatGPTClient struct {
	URL     string
	ApiKey  string
	Model   string
	Context []GPTMesssage
}

func (client *ChatGPTClient) combineResponseChoices(choices []Choice) string {

	var output strings.Builder
	for _, choice := range choices {
		output.WriteString(choice.Message.Content)
		client.Context = append(client.Context, choice.Message)
	}
	return output.String()

}

func (client *ChatGPTClient) defaultRequest() ChatGPTRequest {

	var request ChatGPTRequest
	request.Model = client.Model
	request.Message = make([]GPTMesssage, 0)
	systemMessage := GPTMesssage{
		Role:    "system",
		Content: "Answer the question if possible and continue the conversation",
	}
	request.Message = append(request.Message, systemMessage)
	return request

}

func (client *ChatGPTClient) ChatCompletion(query string) (string, error) {

	chatGPTRequest := client.defaultRequest()
	userMessage := GPTMesssage{
		Role:    "user",
		Content: query,
	}
	chatGPTRequest.Message = append(chatGPTRequest.Message, client.Context...)
	chatGPTRequest.Message = append(chatGPTRequest.Message, userMessage)

	postBody, err := json.Marshal(chatGPTRequest)
	if err != nil {
		log.Printf("couldn't convert request to json %s", err.Error())
		return "", err
	}
	requestBody := bytes.NewBuffer(postBody)

	request, _ := http.NewRequest("POST", client.URL, requestBody)
	authToken := fmt.Sprintf("Bearer %s", client.ApiKey)
	request.Header.Add("Authorization", authToken)
	request.Header.Add("Content-Type", "application/json")

	httpClient := &http.Client{}

	response, err := httpClient.Do(request)
	if err != nil {
		log.Printf("error fetching ChatGPT response %s", err)
		return "", err
	}
	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		log.Printf("error fetching ChatGPT response body %s", err)
		return "", err
	}

	if response.StatusCode != 200 {
		errMessage := "ChatGPT api call returned non 200 statuscode"
		log.Printf("%s %d %s", errMessage, response.StatusCode, string(body))
		return "", fmt.Errorf(errMessage)
	}

	log.Printf("ChatGPT respones for query %s: %s", query, string(body))

	var chatGPTResponse ChatGPTResponse
	err = json.Unmarshal(body, &chatGPTResponse)
	if err != nil {
		log.Printf("error converting ChatGPT response json %s", err)
		return "", err
	}

	reponseText := client.combineResponseChoices(chatGPTResponse.Choices)
	return reponseText, nil

}
