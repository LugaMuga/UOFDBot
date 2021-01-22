package main

type ChatCallback struct {
	Id              int64
	ChatId          int64
	Text            string
	CreateTimestamp int64
}

func (chatCallback *ChatCallback) fill(chatId int64, text string, timestamp int64) {
	chatCallback.ChatId = chatId
	chatCallback.Text = text
	chatCallback.CreateTimestamp = timestamp
}
