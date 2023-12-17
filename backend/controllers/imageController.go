package controllers

import (
	"fmt"
	"log"
	"mime"
	"path/filepath"
	"sync"

	"github.com/google/uuid"
)

type ImageMetaData struct {
	Id       string `json:"Id"`
	FileName string `json:"fileName"`
	FileType string `json:"fileType"`
}

type ImageController struct {
	Images         map[string][]byte
	ImageMetaDatas map[string]ImageMetaData
}

var imageControllerLock = &sync.Mutex{}

var imageController *ImageController

func GetImageControllerInstance() *ImageController {
	if imageController == nil {
		imageControllerLock.Lock()
		defer imageControllerLock.Unlock()
		if imageController == nil {
			imageController = &ImageController{
				Images:         make(map[string][]byte),
				ImageMetaDatas: make(map[string]ImageMetaData),
			}
		}
	}
	return imageController
}

func (ic *ImageController) AddImage(fileName string, image []byte) (string, error) {

	uuid := uuid.NewString()
	log.Printf("Adding image with id: %s", uuid)
	ic.Images[uuid] = image
	metadata := ImageMetaData{
		Id:       uuid,
		FileName: fileName,
		FileType: mime.TypeByExtension(filepath.Ext(fileName)),
	}
	ic.ImageMetaDatas[uuid] = metadata
	log.Printf("image metadata: %+v", metadata)
	log.Printf("No of images: %d", len(ic.Images))
	return uuid, nil

}

func (ic *ImageController) GetImageAndMetaData(id string) (ImageMetaData, []byte, error) {

	if image, ok := ic.Images[id]; !ok {
		errMessage := fmt.Sprintf("no image with given id %s", id)
		log.Printf(errMessage)
		return ImageMetaData{}, nil, fmt.Errorf(errMessage)
	} else {
		metaData := ic.ImageMetaDatas[id]
		return metaData, image, nil
	}

}
