package controllers

import (
	"alle/chat"
	"log"
	"sync"
)

type ChatController struct {
	Chats []*chat.Chat
}

var chatControllerLock = &sync.Mutex{}

var chatController *ChatController

func GetChatControllerInstance() *ChatController {
	if chatController == nil {
		chatControllerLock.Lock()
		defer chatControllerLock.Unlock()
		if chatController == nil {
			chatController = &ChatController{Chats: make([]*chat.Chat, 0)}
		}
	}
	return chatController
}

func (cc *ChatController) NewChat(chat chat.Chat) error {

	cc.Chats = append(cc.Chats, &chat)
	log.Printf("x: %d", len(cc.Chats))
	return nil

}

func (cc *ChatController) AllChat() ([]*chat.Chat, error) {
	return cc.Chats, nil
}
