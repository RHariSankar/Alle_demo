package main

import (
	"alle/client/azure"
	chatgpt "alle/client/chatGPT"
	"alle/routes"
	"log"
	"net/http"
)

func main() {

	azureCLUClient := azure.AzureCLUClient{
		APIMSubscriptionKey: AZURE_OCP_APIM_SUBSCRIPTION_KEY,
		APIMRequestId:       AZURE_APIM_REQUEST_ID,
		URL:                 AZURE_CLU_URL,
		Kind:                AZURE_CLU_KIND,
		ProjectName:         AZURE_CLU_PROJECT_NAME,
		DeploymentName:      AZURE_CLU_DEPLOYMENT_NAME,
		StringIndexType:     AZURE_CLU_STRING_INDEX_TYPE,
		Threshold:           AZURE_CLU_INTENT_THRESHOLD,
		RequiredIntent:      AZURE_REQUIRED_INTENT,
	}

	chatgptClient := chatgpt.ChatGPTClient{
		ApiKey:  CHATGPT_API_KEY,
		Model:   CHATGPT_MODEL,
		URL:     CHATGPT_URL,
		Context: make([]chatgpt.GPTMesssage, 0),
	}

	router := routes.Router(chatgptClient, azureCLUClient)

	log.Println("Server running at port 8080")
	err := http.ListenAndServe(":8080", routes.CorsMiddleware(router))
	if err != nil {
		log.Fatalf("There's an error with the server, %s", err)
	}

}
