package handlers

import (
	"alle/controllers"
	"alle/modals"
	"encoding/json"
	"io"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

type ImageHandler struct {
	ImageController *controllers.ImageController
	ChatController  *controllers.ChatController
}

func (ih *ImageHandler) AddImage(writer http.ResponseWriter, request *http.Request) {

	err := request.ParseMultipartForm(10 << 20) // 10MB max
	if err != nil {
		http.Error(writer, "Error parsing form", http.StatusBadRequest)
		return
	}

	tags := request.Form["tag"]
	if len(tags) == 0 {
		http.Error(writer, "No tags present", http.StatusBadRequest)
		return
	}

	image, handler, err := request.FormFile("image")
	if err != nil {
		http.Error(writer, "Error retrieving file from form", http.StatusBadRequest)
		return
	}
	defer image.Close()
	imageData, err := io.ReadAll(image)
	if err != nil {
		http.Error(writer, "Error reading file", http.StatusInternalServerError)
		return
	}

	id, imageTags, _ := ih.ImageController.AddImage(handler.Filename, imageData, tags)
	imageMessage := modals.ImageChat{
		Role:     "user",
		ImageId:  id,
		DateTime: time.Now().Format(time.RFC3339Nano),
		Tags:     imageTags,
	}
	imageMessage.Type = imageMessage.GetType()
	imageResponse := make([]modals.Chat, 0)
	imageResponse = append(imageResponse, &imageMessage)
	writer.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(writer).Encode(&imageResponse)
	if err != nil {
		http.Error(writer, "Couldn't encode image response", http.StatusInternalServerError)
		return
	}
	ih.ChatController.NewChat(&imageMessage)
}

func (ih *ImageHandler) GetImage(writer http.ResponseWriter, request *http.Request) {

	vars := mux.Vars(request)
	id := vars["id"]

	metaData, image, err := ih.ImageController.GetImageAndMetaData(id)

	if err != nil {
		http.Error(writer, "Image not found", http.StatusNotFound)
		return
	}
	writer.Header().Set("Content-Type", metaData.FileType)
	writer.Write(image)

}

func (ih *ImageHandler) GetImagesByTag(writer http.ResponseWriter, request *http.Request) {

	queryParams := request.URL.Query()
	tags, exist := queryParams["tag"]
	if !exist {
		http.Error(writer, "Missing query parameter 'tags'", http.StatusBadRequest)
		return
	}
	images, _ := ih.ImageController.GetImagesByTags(tags)

	writer.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(writer).Encode(images)
	if err != nil {
		http.Error(writer, "couldn't convert response to json", http.StatusInternalServerError)
		return
	}

}
