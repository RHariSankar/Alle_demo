package controllers

import (
	"alle/modals"
	"fmt"
	"log"
	"mime"
	"path/filepath"
	"sync"
	"time"

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
	ImageTagsMap   map[string]map[string]bool
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
				ImageTagsMap:   make(map[string]map[string]bool),
			}
		}
	}
	return imageController
}

func (ic *ImageController) AddTagToImage(tags []string, imageId string) error {

	for _, tag := range tags {

		if currImageIdMap, ok := ic.ImageTagsMap[tag]; ok {
			currImageIdMap[imageId] = true
		} else {
			ic.ImageTagsMap[tag] = make(map[string]bool)
			ic.ImageTagsMap[tag][imageId] = true
		}

	}
	return nil
}

func (ic *ImageController) AddImage(fileName string, image []byte, tags []string) (string, []string, error) {

	imageId := uuid.NewString()
	log.Printf("Adding image with id: %s", imageId)
	ic.Images[imageId] = image
	metadata := ImageMetaData{
		Id:       imageId,
		FileName: fileName,
		FileType: mime.TypeByExtension(filepath.Ext(fileName)),
	}
	ic.ImageMetaDatas[imageId] = metadata
	ic.AddTagToImage(tags, imageId)
	log.Printf("image metadata: %+v", metadata)
	return imageId, tags, nil

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

func (ic *ImageController) GetImagesByTag(tag string) ([]modals.Chat, error) {

	if imageIds, ok := ic.ImageTagsMap[tag]; ok {
		keys := make([]string, 0, len(imageIds))
		for k := range imageIds {
			keys = append(keys, k)
		}
		images := make([]modals.Chat, 0)
		for _, id := range keys {
			image := modals.ImageChat{
				Type:     "image",
				Role:     "system",
				ImageId:  id,
				DateTime: time.Now().Format(time.RFC3339),
				Tags:     []string{tag},
			}
			images = append(images, &image)
		}
		return images, nil
	} else {
		message := modals.TextChat{
			Type:     "chat",
			Role:     "system",
			Text:     "No image with tag " + tag,
			DateTime: time.Now().Format(time.RFC3339),
		}
		return []modals.Chat{&message}, nil
	}

}

func (ic *ImageController) GetImagesByTags(tags []string) ([]modals.Chat, error) {

	images := make([]modals.Chat, 0)
	for _, tag := range tags {
		imagesForTag, _ := ic.GetImagesByTag(tag)
		images = append(images, imagesForTag...)
	}
	return images, nil

}
