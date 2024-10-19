package services

import (
	"github.com/LugaMuga/UOFDBot/internal/bot"
	"github.com/LugaMuga/UOFDBot/internal/config"
	"github.com/LugaMuga/UOFDBot/internal/dao"
	"github.com/LugaMuga/UOFDBot/internal/locale"
	"github.com/LugaMuga/UOFDBot/internal/models"
	"github.com/LugaMuga/UOFDBot/internal/utils"
	tgbotapi "github.com/Syfaro/telegram-bot-api"
	"strconv"
	"strings"
)

const ResetPidorPoll = `RESET_PIDOR`
const ResetHeroPoll = `RESET_HERO`
const ResetPollAgreedOption = `yes`
const ResetPollDisagreedOption = `no`

func ResetApproval(chatId int64, resetPollName string) {
	resetPoll := NewSimplePoll(resetPollName)
	msg := tgbotapi.NewMessage(chatId, locale.Loc(locale.DefaultLang, `a_u_sure`, config.Config.BotResetMinPercentage))
	chatCallback := new(models.ChatCallback)
	chatCallback.Fill(chatId, "", utils.NowUnix())
	msg.ReplyMarkup = buildResetPollMarkup(resetPoll, chatCallback)
	_, _ = bot.Bot.Send(msg)
}

func ResetPoll(update tgbotapi.Update) {
	message := update.CallbackQuery.Message
	params := strings.Split(update.CallbackQuery.Data, CallbackQueryParamDelimiter)
	n, _ := strconv.Atoi(params[3])
	callbackDataDb := dao.GetChatCallbackById(n)
	poll := ParseSimplePollCallbackQuery(update.CallbackQuery)
	if poll.name != ResetPidorPoll && poll.name != ResetHeroPoll {
		return
	}
	activeChatUsers := dao.GetEnabledChatUsersByChatId(message.Chat.ID)
	var activeUsernames []string
	for _, activeChatUser := range activeChatUsers {
		checkAndUpdateUserIfNeeded(&activeChatUser, update.CallbackQuery.From, message.Chat.ID)
		activeUsernames = append(activeUsernames, utils.FormatChatUserName(activeChatUser))
	}

	poll.applySelectedOption(update.CallbackQuery.From, ResetPollAgreedOption, ResetPollDisagreedOption)
	poll.improveVotedUserArrays(activeUsernames)
	agreedPercentage := calcPercentage(len(poll.agreedUsers), len(activeUsernames))
	disagreedPercentage := calcPercentage(len(poll.disagreedUsers), len(activeUsernames))
	poll.updateButtonsText(agreedPercentage, disagreedPercentage)

	if agreedPercentage >= config.Config.BotResetMinPercentage {
		resetByPollName(message.Chat.ID, message.MessageID, update.CallbackQuery.ID, poll.name, callbackDataDb.Id)
		return
	} else if disagreedPercentage >= config.Config.BotResetMinPercentage {
		resetByPollName(message.Chat.ID, message.MessageID, update.CallbackQuery.ID, `nil`, callbackDataDb.Id)
		return
	}
	editedMarkup := buildResetPollMarkup(poll, callbackDataDb)
	editedMsg := tgbotapi.NewEditMessageReplyMarkup(message.Chat.ID, message.MessageID, editedMarkup)
	_, _ = bot.Bot.Send(editedMsg)
}

func checkAndUpdateUserIfNeeded(activeChatUser *models.ChatUser, from *tgbotapi.User, chatId int64) {
	if activeChatUser.UserId != from.ID ||
		(activeChatUser.Username == from.UserName &&
			activeChatUser.UserFirstName == from.FirstName &&
			activeChatUser.UserLastName == from.LastName) {
		return
	}
	activeChatUser.Username = from.UserName
	activeChatUser.UserFirstName = from.FirstName
	activeChatUser.UserLastName = from.LastName
	activeChatUser.ChatId = chatId
	dao.UpdateChatUserUsernameFirstAndLastName(activeChatUser)
}

func calcPercentage(part int, full int) int {
	if full == 0 {
		return 0
	}
	return int(float32(part) / float32(full) * 100)
}

func resetByPollName(chatId int64, msgId int, callbackId string, pollName string, chatCallbackId int64) {
	switch pollName {
	case ResetPidorPoll:
		resetPidor(chatId)
		break
	case ResetHeroPoll:
		resetHero(chatId)
		break
	default:
		msg := locale.Loc(locale.DefaultLang, `stat_not_reset`)
		bot.Bot.AnswerCallbackQuery(tgbotapi.NewCallback(callbackId, msg))
		bot.SendMessage(chatId, msg)
	}
	deleteKeyBoardMsg := tgbotapi.NewDeleteMessage(chatId, msgId)
	dao.DeleteChatCallback(chatCallbackId, chatId)
	_, _ = bot.Bot.Send(deleteKeyBoardMsg)
}

func buildResetPollMarkup(resetPoll *SimplePoll, chatCallback *models.ChatCallback) tgbotapi.InlineKeyboardMarkup {
	if chatCallback.Id == 0 {
		id := dao.SaveOrUpdateChatCallback(*chatCallback)
		chatCallback.Id = id
	}
	agreedButtonData := buildSimplePollButtonData(resetPoll, ResetPollAgreedOption, chatCallback)
	disagreedButtonData := buildSimplePollButtonData(resetPoll, ResetPollDisagreedOption, chatCallback)
	dao.SaveOrUpdateChatCallback(*chatCallback)
	return tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(resetPoll.agreedText, agreedButtonData),
			tgbotapi.NewInlineKeyboardButtonData(resetPoll.disagreedText, disagreedButtonData),
		))
}
