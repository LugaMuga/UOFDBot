package main

import (
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
	msg := tgbotapi.NewMessage(chatId, loc(defaultLang, `a_u_sure`, config.BotResetMinPercentage))
	chatCallback := new(ChatCallback)
	chatCallback.fill(chatId, "", nowUnix())
	msg.ReplyMarkup = buildResetPollMarkup(resetPoll, chatCallback)
	_, _ = bot.Send(msg)
}

func ResetPoll(update tgbotapi.Update) {
	message := update.CallbackQuery.Message
	params := strings.Split(update.CallbackQuery.Data, CallbackQueryParamDelimiter)
	n, _ := strconv.Atoi(params[3])
	callbackDataDb := getChatCallbackById(n)
	poll := ParseSimplePollCallbackQuery(update.CallbackQuery)
	if poll.name != ResetPidorPoll && poll.name != ResetHeroPoll {
		return
	}
	activeChatUsers := getEnabledChatUsersByChatId(message.Chat.ID)
	var activeUsernames []string
	for _, activeChatUser := range activeChatUsers {
		checkAndUpdateUserIfNeeded(&activeChatUser, update.CallbackQuery.From, message.Chat.ID)
		activeUsernames = append(activeUsernames, FormatChatUserName(activeChatUser))
	}

	poll.applySelectedOption(update.CallbackQuery.From, ResetPollAgreedOption, ResetPollDisagreedOption)
	poll.improveVotedUserArrays(activeUsernames)
	agreedPercentage := calcPercentage(len(poll.agreedUsers), len(activeUsernames))
	disagreedPercentage := calcPercentage(len(poll.disagreedUsers), len(activeUsernames))
	poll.updateButtonsText(agreedPercentage, disagreedPercentage)

	if agreedPercentage >= config.BotResetMinPercentage {
		resetByPollName(message.Chat.ID, message.MessageID, update.CallbackQuery.ID, poll.name, callbackDataDb.Id)
		return
	} else if disagreedPercentage >= config.BotResetMinPercentage {
		resetByPollName(message.Chat.ID, message.MessageID, update.CallbackQuery.ID, `nil`, callbackDataDb.Id)
		return
	}
	editedMarkup := buildResetPollMarkup(poll, callbackDataDb)
	editedMsg := tgbotapi.NewEditMessageReplyMarkup(message.Chat.ID, message.MessageID, editedMarkup)
	_, _ = bot.Send(editedMsg)
}

func checkAndUpdateUserIfNeeded(activeChatUser *ChatUser, from *tgbotapi.User, chatId int64) {
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
	UpdateChatUserUsernameFirstAndLastName(activeChatUser)
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
		msg := loc(defaultLang, `stat_not_reset`)
		bot.AnswerCallbackQuery(tgbotapi.NewCallback(callbackId, msg))
		SendMessage(chatId, msg)
	}
	deleteKeyBoardMsg := tgbotapi.NewDeleteMessage(chatId, msgId)
	DeleteChatCallback(chatCallbackId, chatId)
	_, _ = bot.Send(deleteKeyBoardMsg)
}

func buildResetPollMarkup(resetPoll *SimplePoll, chatCallback *ChatCallback) tgbotapi.InlineKeyboardMarkup {
	if chatCallback.Id == 0 {
		id := SaveOrUpdateChatCallback(*chatCallback)
		chatCallback.Id = id
	}
	agreedButtonData := buildSimplePollButtonData(resetPoll, ResetPollAgreedOption, chatCallback)
	disagreedButtonData := buildSimplePollButtonData(resetPoll, ResetPollDisagreedOption, chatCallback)
	SaveOrUpdateChatCallback(*chatCallback)
	return tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(resetPoll.agreedText, agreedButtonData),
			tgbotapi.NewInlineKeyboardButtonData(resetPoll.disagreedText, disagreedButtonData),
		))
}
