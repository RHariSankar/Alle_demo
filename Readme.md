# Alle Assignment

## Assumption
- Images are always stored against some tags
- Image retreival is always based on the stored tag
    - Azure CLU model is trained with a few generic prompt for image retreival
    - It is configured to obtain the entity (tag for image) from the natural language query
    - Returned entity is used for image retrieval

## Backend

### Description
- Go based server using Gorilla/MUX
- All data are stored in in-memory data structures
- Calls 2 external services 
  1. Azure Conversational Language Understanding - To enable language understanding use case (image retreival)
  2. ChatGPT - To enable general conversations

### Build
Replace the constants in `conf.go` with actual values

```go mod download```

```go build . && ./alle```

## Frontend

### Description
- Simple chat interface using React
- Interacts with server through API

### Build
```npm install && npm start```

## Reference
[Azure Conversational Language Understanding](https://learn.microsoft.com/en-us/azure/ai-services/language-service/conversational-language-understanding/overview)

[ChatGPT API](https://platform.openai.com/docs/introduction)

[Gorilla/MUX](https://pkg.go.dev/github.com/gorilla/mux)