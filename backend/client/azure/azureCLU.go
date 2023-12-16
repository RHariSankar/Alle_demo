package azure

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)

type ConversationItem struct {
	Id            string `json:"id,omitempty"`
	Text          string `json:"text"`
	Modality      string `json:"modality,omitempty"`
	ParticipantId string `json:"participantId,omitempty"`
}

type AnalysisInput struct {
	ConversationItem ConversationItem `json:"conversationItem"`
}

type Parameters struct {
	ProjectName     string `json:"projectName"`
	Verbose         bool   `json:"verbose,omitempty"`
	DeploymentName  string `json:"deploymentName"`
	StringIndexType string `json:"stringIndexType,omitempty"`
}

type AzureCLURequest struct {
	Kind          string        `json:"kind"`
	AnalysisInput AnalysisInput `json:"analysisInput"`
	Parameters    Parameters    `json:"parameters"`
}

type Intent struct {
	Category        string  `json:"category"`
	ConfidenceScore float32 `json:"confidenceScore"`
}

type Entity struct {
	Category        string  `json:"category"`
	Text            string  `json:"text"`
	ConfidenceScore float32 `json:"confidenceScore"`
}

type Prediction struct {
	TopIntent string   `json:"topIntent"`
	Intents   []Intent `json:"intents"`
	Entities  []Entity `json:"entities"`
}

type Result struct {
	Query      string     `json:"query"`
	Prediction Prediction `json:"prediction"`
}

type AzureCLUResponse struct {
	Kind   string `json:"kind"`
	Result Result `json:"result"`
}

type AzureCLU struct {
	APIMSubscriptionKey string
	APIMRequestId       string
	URL                 string
	Kind                string
	ProjectName         string
	DeploymentName      string
	StringIndexType     string
	IntentThreshold     float32
	RequiredIntent      string
}

func (clu *AzureCLU) defaultRequest() AzureCLURequest {
	azureCLURequest := AzureCLURequest{}
	azureCLURequest.Kind = clu.Kind
	azureCLURequest.AnalysisInput.ConversationItem.Id = "1"
	azureCLURequest.AnalysisInput.ConversationItem.ParticipantId = "1"
	azureCLURequest.Parameters.ProjectName = clu.ProjectName
	azureCLURequest.Parameters.DeploymentName = clu.DeploymentName
	azureCLURequest.Parameters.Verbose = true
	azureCLURequest.Parameters.StringIndexType = clu.StringIndexType

	return azureCLURequest
}

func (clu *AzureCLU) getIntentAndEntity(azureCLUResponse AzureCLUResponse) (bool, string) {

	if azureCLUResponse.Result.Prediction.TopIntent != clu.RequiredIntent {
		return false, ""
	}

	for _, intent := range azureCLUResponse.Result.Prediction.Intents {
		if intent.Category == clu.RequiredIntent && intent.ConfidenceScore < clu.IntentThreshold {
			return false, ""
		}
	}

	if len(azureCLUResponse.Result.Prediction.Entities) > 0 {
		return true, azureCLUResponse.Result.Prediction.Entities[0].Text
	} else {
		return true, "all"
	}

}

func (clu *AzureCLU) GetIntentAndEntity(query string) (bool, string, error) {

	azureCLUrequest := clu.defaultRequest()
	azureCLUrequest.AnalysisInput.ConversationItem.Text = query

	postBody, err := json.Marshal(azureCLUrequest)
	if err != nil {
		log.Printf("couldn't convert request to json %s", err.Error())
		return false, "", err
	}
	requestBody := bytes.NewBuffer(postBody)

	request, _ := http.NewRequest("POST", clu.URL, requestBody)
	request.Header.Add("Ocp-Apim-Subscription-Key", clu.APIMSubscriptionKey)
	request.Header.Add("Apim-Request-Id", clu.APIMRequestId)
	request.Header.Add("Content-Type", "application/json")

	client := &http.Client{}

	response, err := client.Do(request)
	if err != nil {
		log.Printf("error fetching Azure CLU response %s", err)
		return false, "", err
	}
	defer response.Body.Close()

	if response.StatusCode != 200 {
		log.Printf("ChatGPT api call returned non 200 statuscode %s", err)
		return false, "", fmt.Errorf("ChatGPT api call returned non 200 statuscode")
	}

	body, err := io.ReadAll(response.Body)
	if err != nil {
		log.Printf("error fetching Azure CLU response body %s", err)
		return false, "", err
	}

	var azureCLUResponse AzureCLUResponse
	log.Printf("AzureCLU Response: %s", string(body))
	err = json.Unmarshal(body, &azureCLUResponse)
	if err != nil {
		log.Printf("error converting Azure CLU response json %s", err)
		return false, "", err
	}

	isIntent, entity := clu.getIntentAndEntity(azureCLUResponse)
	return isIntent, entity, nil

}
