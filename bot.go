package main

import (
	tgbotapi "github.com/Syfaro/telegram-bot-api"
	"log"
	"net/http"
	"strings"
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
	if update.Message != nil && update.Message.LeftChatMember != nil {
		applyLeftChatMember(update)
	} else if update.Message != nil {
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

	log.Printf("[%s] %s", FormatUserNameFromApi(update.Message.From), update.Message.Text)

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
		delete(message.Chat.ID, message.From)
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
		ResetApproval(message.Chat.ID, ResetPidorPoll)
		break
	case `/hero`, `/hero@` + config.BotName:
		hero(message.Chat.ID)
		break
	case `/herolist`, `/herolist@` + config.BotName:
		heroList(message.Chat.ID)
		break
	case `/resethero`, `/resethero@` + config.BotName:
		ResetApproval(message.Chat.ID, ResetHeroPoll)
		break
	case `/update`, `/update@` + config.BotName:
		updateUsers(message.Chat.ID)
		break
	default:
		break
	}
}

func applyCallbackQuery(update tgbotapi.Update) {
	if update.CallbackQuery == nil {
		return
	}
	callbackQueryType := getCallbackQueryType(update.CallbackQuery.Data)
	switch callbackQueryType {
	case string(SimplePollType):
		ResetPoll(update)
		return
	}
}

func getCallbackQueryType(data string) string {
	i := strings.Index(data, CallbackQueryParamDelimiter)
	return data[0:i]
}

func applyLeftChatMember(update tgbotapi.Update) {
	message := update.Message
	log.Printf("[%s] left chat", FormatUserNameFromApi(update.Message.LeftChatMember))
	delete(message.Chat.ID, message.LeftChatMember)
}

func SendMessage(chatID int64, text string) {
	msg := tgbotapi.NewMessage(chatID, text)
	_, _ = bot.Send(msg)
}
