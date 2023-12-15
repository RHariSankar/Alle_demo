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

type AzureCLU struct{}

func defaultRequest() AzureCLURequest {
	azureCLURequest := AzureCLURequest{}
	azureCLURequest.Kind = AZURE_CLU_KIND
	azureCLURequest.AnalysisInput.ConversationItem.Id = "1"
	azureCLURequest.AnalysisInput.ConversationItem.ParticipantId = "1"
	azureCLURequest.Parameters.ProjectName = AZURE_CLU_PROJECT_NAME
	azureCLURequest.Parameters.DeploymentName = AZURE_CLU_DEPLOYMENT_NAME
	azureCLURequest.Parameters.Verbose = true
	azureCLURequest.Parameters.StringIndexType = AZURE_CLU_STRING_INDEX_TYPE

	return azureCLURequest
}

func getIntentAndEntity(azureCLUResponse AzureCLUResponse) (bool, string) {

	if azureCLUResponse.Result.Prediction.TopIntent != ALLE_REQUIRED_INTENT {
		return false, ""
	}

	for _, intent := range azureCLUResponse.Result.Prediction.Intents {
		if (intent.Category == ALLE_REQUIRED_INTENT && intent.ConfidenceScore < AZURE_CLU_INTENT_THRESHOLD) {
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

	azureCLUrequest := defaultRequest()
	azureCLUrequest.AnalysisInput.ConversationItem.Text = query

	postBody, err := json.Marshal(azureCLUrequest)
	if err != nil {
		return false, "", err
	}
	requestBody := bytes.NewBuffer(postBody)

	request, _ := http.NewRequest("POST", AZURE_CLU_URL, requestBody)
	request.Header.Add(AZURE_OCP_APIM_SUBSCRIPTION_KEY_HEADER, AZURE_OCP_APIM_SUBSCRIPTION_KEY)
	request.Header.Add(AZURE_APIM_REQUEST_ID_HEADER, AZURE_APIM_REQUEST_ID)
	request.Header.Add("Content-Type", "application/json")

	client := &http.Client{}

	response, err := client.Do(request)
	if err != nil {
		log.Printf("error fetching Azure CLU response %s", err)
		return false, "", err
	}
	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		log.Printf("error fetching Azure CLU response body %s", err)
		return false, "", err
	}

	var azureCLUResponse AzureCLUResponse
	fmt.Println(string(body))
	err = json.Unmarshal(body, &azureCLUResponse)
	if err != nil {
		log.Printf("error converting Azure CLU response json %s", err)
		return false, "", err
	}

	isIntent, entity := getIntentAndEntity(azureCLUResponse)
	return isIntent, entity, nil

}
