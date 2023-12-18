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

type AzureCLUClient struct {
	APIMSubscriptionKey string
	APIMRequestId       string
	URL                 string
	Kind                string
	ProjectName         string
	DeploymentName      string
	StringIndexType     string
	Threshold           float32
	RequiredIntent      string
}

type AzureCLUNoEnitityError struct{}

func (e *AzureCLUNoEnitityError) Error() string {
	return "Couldn't find entities in query"
}

func (clu *AzureCLUClient) defaultRequest() AzureCLURequest {
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

func (clu *AzureCLUClient) getIntentAndEntity(azureCLUResponse AzureCLUResponse) (bool, []string) {

	if azureCLUResponse.Result.Prediction.TopIntent != clu.RequiredIntent {
		return false, nil
	}

	for _, intent := range azureCLUResponse.Result.Prediction.Intents {
		if intent.Category == clu.RequiredIntent && intent.ConfidenceScore < clu.Threshold {
			return false, nil
		}
	}

	if len(azureCLUResponse.Result.Prediction.Entities) > 0 {
		entities := make([]string, 0)
		for _, entity := range azureCLUResponse.Result.Prediction.Entities {
			if entity.ConfidenceScore >= clu.Threshold {
				entities = append(entities, entity.Text)
			}
		}
		return true, entities
	} else {
		return true, nil
	}

}

func (clu *AzureCLUClient) GetIntentAndEntity(query string) (bool, []string, error) {

	azureCLUrequest := clu.defaultRequest()
	azureCLUrequest.AnalysisInput.ConversationItem.Text = query

	postBody, err := json.Marshal(azureCLUrequest)
	if err != nil {
		log.Printf("couldn't convert request to json %s", err.Error())
		return false, nil, err
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
		return false, nil, err
	}
	defer response.Body.Close()

	if response.StatusCode != 200 {
		log.Printf("AzureCLU api call returned non 200 statuscode %s", err)
		return false, nil, fmt.Errorf("AzureCLU api call returned non 200 statuscode")
	}

	body, err := io.ReadAll(response.Body)
	if err != nil {
		log.Printf("error fetching Azure CLU response body %s", err)
		return false, nil, err
	}

	var azureCLUResponse AzureCLUResponse
	log.Printf("AzureCLU Response for query %s: %s", query, string(body))
	err = json.Unmarshal(body, &azureCLUResponse)
	if err != nil {
		log.Printf("error converting Azure CLU response json %s", err)
		return false, nil, err
	}

	isIntent, entity := clu.getIntentAndEntity(azureCLUResponse)
	if isIntent && entity == nil {
		return false, nil, &AzureCLUNoEnitityError{}
	}
	return isIntent, entity, nil

}
