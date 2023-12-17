package handlers

import (
	"alle/chat"
	"alle/controllers"
	"encoding/json"
	"io"
	"log"
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

	log.Printf("handler %+v", handler.Filename)

	id, _ := ih.ImageController.AddImage(handler.Filename, imageData)
	imageMessage := chat.ImageChat{
		Role:     "user",
		ImageId:  id,
		DateTime: time.Now().Format(time.RFC3339Nano),
	}
	imageMessage.Type = imageMessage.GetType()
	writer.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(writer).Encode(&imageMessage)
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
