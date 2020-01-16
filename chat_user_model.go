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

func (chatUser *ChatUser) fill(chatId int64, user *tgbotapi.User) {
	chatUser.chatId = chatId
	chatUser.userId = user.ID
	chatUser.username = user.UserName
	chatUser.userFirstName = user.FirstName
	chatUser.userLastName = user.LastName
}
