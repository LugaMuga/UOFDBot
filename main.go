package main

import (
	"log"
	"net/http"
	"strings"

	"github.com/LugaMuga/UOFDBot/internal/bot"
	"github.com/LugaMuga/UOFDBot/internal/config"
	"github.com/LugaMuga/UOFDBot/internal/db"
	"github.com/LugaMuga/UOFDBot/internal/locale"
	"github.com/LugaMuga/UOFDBot/internal/services"
	"github.com/LugaMuga/UOFDBot/internal/utils"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
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
	webhookConfig, err := tgbotapi.NewWebhookWithCert("https://"+config.Config.WebhookHost+":"+config.Config.WebhookPort+"/"+bot.Bot.Token, tgbotapi.FilePath(config.Config.WebhookCertPath))
	if err != nil {
		log.Fatalf("Ошибка создания конфига вебхука: %v", err)
	}
	_, err = bot.Bot.Request(webhookConfig)
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

	updates := bot.Bot.GetUpdatesChan(updateConfig)

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
	if update.Message == nil || !update.Message.IsCommand() {
		return
	}

	log.Printf("[%s] %s", utils.FormatUserNameFromApi(update.Message.From), update.Message.Text)

	message := update.Message
	switch message.Command() {
	case `start`:
		msg := locale.Loc(locale.DefaultLang, `hello`)
		bot.SendMessage(update.Message.Chat.ID, msg)

	case `register`:
		services.Register(*message)

	case `delete`:
		services.Delete(message.Chat.ID, message.From)

	case `run`:
		services.Run(message.Chat.ID)

	case `list`:
		services.List(message.Chat.ID)

	case `pidor`:
		services.Pidor(message.Chat.ID)

	case `pidorlist`:
		services.PidorList(message.Chat.ID)

	case `resetpidor`:
		services.ResetApproval(message.Chat.ID, services.ResetPidorPoll)

	case `hero`:
		services.Hero(message.Chat.ID)

	case `herolist`:
		services.HeroList(message.Chat.ID)

	case `resethero`:
		services.ResetApproval(message.Chat.ID, services.ResetHeroPoll)

	case `update`:
		services.UpdateUsers(message.Chat.ID)

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
