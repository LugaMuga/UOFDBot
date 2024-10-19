package models

import tgbotapi "github.com/Syfaro/telegram-bot-api"

type ChatUser struct {
	Id                 int
	ChatId             int64
	UserId             int
	Username           string
	UserFirstName      string
	UserLastName       string
	Enabled            bool
	PidorScore         int
	PidorLastTimestamp int64
	HeroScore          int
	HeroLastTimestamp  int64
}

func (chatUser *ChatUser) Fill(chatId int64, user *tgbotapi.User) {
	chatUser.ChatId = chatId
	chatUser.UserId = user.ID
	chatUser.Username = user.UserName
	chatUser.UserFirstName = user.FirstName
	chatUser.UserLastName = user.LastName
}
