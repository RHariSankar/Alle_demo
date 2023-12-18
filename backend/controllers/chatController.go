package controllers

import (
	"alle/modals"
	"log"
	"sync"
)

type ChatController struct {
	Chats []*modals.Chat
}

var chatControllerLock = &sync.Mutex{}

var chatController *ChatController

func GetChatControllerInstance() *ChatController {
	if chatController == nil {
		chatControllerLock.Lock()
		defer chatControllerLock.Unlock()
		if chatController == nil {
			chatController = &ChatController{Chats: make([]*modals.Chat, 0)}
		}
	}
	return chatController
}

func (cc *ChatController) NewChat(chat modals.Chat) error {

	cc.Chats = append(cc.Chats, &chat)
	log.Printf("x: %d", len(cc.Chats))
	return nil

}

func (cc *ChatController) AllChat() ([]*modals.Chat, error) {
	return cc.Chats, nil
}
