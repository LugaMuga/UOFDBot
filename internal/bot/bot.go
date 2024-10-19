package bot

import (
	"github.com/LugaMuga/UOFDBot/internal/config"
	tgbotapi "github.com/Syfaro/telegram-bot-api"
	"log"
)

var Bot *tgbotapi.BotAPI

func RegisterBot() {
	var err error
	Bot, err = tgbotapi.NewBotAPI(config.Config.BotToken)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Authorized on account %s", Bot.Self.UserName)
	Bot.RemoveWebhook()
}

func SendMessage(chatID int64, text string) {
	msg := tgbotapi.NewMessage(chatID, text)
	_, _ = Bot.Send(msg)
}
