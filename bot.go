package main

import (
	tgbotapi "github.com/Syfaro/telegram-bot-api"
	"log"
	"net/http"
)

var bot *tgbotapi.BotAPI

func main() {
	LoadConfig()
	LoadI18N()
	InitDb()
	registerBot()
	if config.ConnectionType == ConnectionTypeWebhook {
		subscribeWithWebhook()
	} else {
		subscribeToUpdatesChan()
	}
	defer CloseDb()
}

func registerBot() {
	var err error
	bot, err = tgbotapi.NewBotAPI(config.BotToken)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Authorized on account %s", bot.Self.UserName)
	bot.RemoveWebhook()
}

func subscribeWithWebhook() {
	_, err := bot.SetWebhook(tgbotapi.NewWebhookWithCert("https://"+config.WebhookHost+":"+config.WebhookPort+"/"+bot.Token, config.WebhookCertPath))
	if err != nil {
		log.Fatal(err)
	}
	info, err := bot.GetWebhookInfo()
	if err != nil {
		log.Fatal(err)
	}
	if info.LastErrorDate != 0 {
		log.Printf("Telegram callback failed: %s", info.LastErrorMessage)
	}
	updates := bot.ListenForWebhook("/" + bot.Token)
	go http.ListenAndServeTLS(config.ServerAddress+":"+config.WebhookPort, config.WebhookCertPath, config.WebhookKeyPath, nil)

	for update := range updates {
		applyUpdate(update)
	}
}

func subscribeToUpdatesChan() {
	updateConfig := tgbotapi.NewUpdate(0)
	updateConfig.Timeout = config.BotTimeout

	updates, _ := bot.GetUpdatesChan(updateConfig)

	for update := range updates {
		applyUpdate(update)
	}
}

func applyUpdate(update tgbotapi.Update) {
	if update.Message != nil {
		applyCommand(update)
	} else if update.CallbackQuery != nil {
		applyCallbackQuery(update)
	}
}

func applyCommand(update tgbotapi.Update) {
	if update.Message == nil {
		return
	}
	if update.Message.Entities == nil ||
		len(*update.Message.Entities) == 0 ||
		(*update.Message.Entities)[0].Type != `bot_command` {
		return
	}

	log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)

	message := update.Message
	switch update.Message.Text {
	case `/start`, `/start@` + config.BotName:
		msg := loc(defaultLang, `hello`)
		SendMessage(update.Message.Chat.ID, msg)
		break
	case `/register`, `/register@` + config.BotName:
		register(*message)
		break
	case `/delete`, `/delete@` + config.BotName:
		delete(*message)
		break
	case `/run`, `/run@` + config.BotName:
		run(message.Chat.ID)
		break
	case `/list`, `/list@` + config.BotName:
		list(message.Chat.ID)
		break
	case `/pidor`, `/pidor@` + config.BotName:
		pidor(message.Chat.ID)
		break
	case `/pidorlist`, `/pidorlist@` + config.BotName:
		pidorList(message.Chat.ID)
		break
	case `/resetpidor`, `/resetpidor@` + config.BotName:
		resetApproval(message.Chat.ID, resetPidorApproval)
		break
	case `/hero`, `/hero@` + config.BotName:
		hero(message.Chat.ID)
		break
	case `/herolist`, `/herolist@` + config.BotName:
		heroList(message.Chat.ID)
		break
	case `/resethero`, `/resethero@` + config.BotName:
		resetApproval(message.Chat.ID, resetHeroApproval)
		break
	default:
		break
	}
}

func applyCallbackQuery(update tgbotapi.Update) {
	if update.CallbackQuery == nil {
		return
	}

	message := update.CallbackQuery.Message
	switch update.CallbackQuery.Data {
	case resetPidorApproval:
		resetPidor(message.Chat.ID)
		break
	case resetHeroApproval:
		resetHero(message.Chat.ID)
		break
	default:
		msg := loc(defaultLang, `stat_not_reset`)
		bot.AnswerCallbackQuery(tgbotapi.NewCallback(update.CallbackQuery.ID, msg))
		SendMessage(message.Chat.ID, msg)
	}
	deleteKeyBoardMsg := tgbotapi.NewDeleteMessage(message.Chat.ID, message.MessageID)
	_, _ = bot.Send(deleteKeyBoardMsg)
}

func SendMessage(chatID int64, text string) {
	msg := tgbotapi.NewMessage(chatID, text)
	_, _ = bot.Send(msg)
}
