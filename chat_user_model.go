package main

import tgbotapi "github.com/Syfaro/telegram-bot-api"

type ChatUser struct {
	id                 int
	chatId             int64
	userId             int
	username           string
	userFirstName      string
	userLastName       string
	enabled            bool
	pidorScore         int
	pidorLastTimestamp int64
	heroScore          int
	heroLastTimestamp  int64
}

func (chatUser *ChatUser) fillFromMessage(message tgbotapi.Message) {
	chatUser.chatId = message.Chat.ID
	chatUser.userId = message.From.ID
	chatUser.username = message.From.UserName
	chatUser.userFirstName = message.From.FirstName
	chatUser.userLastName = message.From.LastName
}
