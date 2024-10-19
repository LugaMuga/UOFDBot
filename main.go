package main

import (
	"github.com/LugaMuga/UOFDBot/internal/bot"
	"github.com/LugaMuga/UOFDBot/internal/config"
	"github.com/LugaMuga/UOFDBot/internal/db"
	"github.com/LugaMuga/UOFDBot/internal/locale"
	"github.com/LugaMuga/UOFDBot/internal/services"
	"github.com/LugaMuga/UOFDBot/internal/utils"
	tgbotapi "github.com/Syfaro/telegram-bot-api"
	"log"
	"net/http"
	"strings"
)

func main() {
	config.LoadConfig()
	locale.LoadI18N()
	db.InitDb()
	bot.RegisterBot()
	if config.Config.ConnectionType == config.ConnectionTypeWebhook {
		subscribeWithWebhook()
	} else {
		subscribeToUpdatesChan()
	}
	defer db.CloseDb()
}

func subscribeWithWebhook() {
	_, err := bot.Bot.SetWebhook(tgbotapi.NewWebhookWithCert("https://"+config.Config.WebhookHost+":"+config.Config.WebhookPort+"/"+bot.Bot.Token, config.Config.WebhookCertPath))
	if err != nil {
		log.Fatal(err)
	}
	info, err := bot.Bot.GetWebhookInfo()
	if err != nil {
		log.Fatal(err)
	}
	if info.LastErrorDate != 0 {
		log.Printf("Telegram callback failed: %s", info.LastErrorMessage)
	}
	updates := bot.Bot.ListenForWebhook("/" + bot.Bot.Token)
	go http.ListenAndServeTLS(config.Config.ServerAddress+":"+config.Config.WebhookPort, config.Config.WebhookCertPath, config.Config.WebhookKeyPath, nil)

	for update := range updates {
		applyUpdate(update)
	}
}

func subscribeToUpdatesChan() {
	updateConfig := tgbotapi.NewUpdate(0)
	updateConfig.Timeout = config.Config.BotTimeout

	updates, _ := bot.Bot.GetUpdatesChan(updateConfig)

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

	log.Printf("[%s] %s", utils.FormatUserNameFromApi(update.Message.From), update.Message.Text)

	message := update.Message
	switch update.Message.Text {
	case `/start`, `/start@` + config.Config.BotName:
		msg := locale.Loc(locale.DefaultLang, `hello`)
		bot.SendMessage(update.Message.Chat.ID, msg)
		break
	case `/register`, `/register@` + config.Config.BotName:
		services.Register(*message)
		break
	case `/delete`, `/delete@` + config.Config.BotName:
		services.Delete(message.Chat.ID, message.From)
		break
	case `/run`, `/run@` + config.Config.BotName:
		services.Run(message.Chat.ID)
		break
	case `/list`, `/list@` + config.Config.BotName:
		services.List(message.Chat.ID)
		break
	case `/pidor`, `/pidor@` + config.Config.BotName:
		services.Pidor(message.Chat.ID)
		break
	case `/pidorlist`, `/pidorlist@` + config.Config.BotName:
		services.PidorList(message.Chat.ID)
		break
	case `/resetpidor`, `/resetpidor@` + config.Config.BotName:
		services.ResetApproval(message.Chat.ID, services.ResetPidorPoll)
		break
	case `/hero`, `/hero@` + config.Config.BotName:
		services.Hero(message.Chat.ID)
		break
	case `/herolist`, `/herolist@` + config.Config.BotName:
		services.HeroList(message.Chat.ID)
		break
	case `/resethero`, `/resethero@` + config.Config.BotName:
		services.ResetApproval(message.Chat.ID, services.ResetHeroPoll)
		break
	case `/update`, `/update@` + config.Config.BotName:
		services.UpdateUsers(message.Chat.ID)
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
	case string(services.SimplePollType):
		services.ResetPoll(update)
		return
	}
}

func getCallbackQueryType(data string) string {
	i := strings.Index(data, services.CallbackQueryParamDelimiter)
	return data[0:i]
}

func applyLeftChatMember(update tgbotapi.Update) {
	message := update.Message
	log.Printf("[%s] left chat", utils.FormatUserNameFromApi(update.Message.LeftChatMember))
	services.Delete(message.Chat.ID, message.LeftChatMember)
}
