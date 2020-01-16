package main

import tgbotapi "github.com/Syfaro/telegram-bot-api"

const ResetPidorPoll = `RESET_PIDOR`
const ResetHeroPoll = `RESET_HERO`
const ResetPollAgreedOption = `yes`
const ResetPollDisagreedOption = `no`

func ResetApproval(chatId int64, resetPollName string) {
	resetPoll := NewSimplePoll(resetPollName)
	msg := tgbotapi.NewMessage(chatId, loc(defaultLang, `a_u_sure`, config.BotResetMinPercentage))
	msg.ReplyMarkup = buildResetPollMarkup(resetPoll)
	_, _ = bot.Send(msg)
}

func ResetPoll(update tgbotapi.Update) {
	message := update.CallbackQuery.Message
	poll := ParseSimplePollCallbackQuery(update.CallbackQuery)
	if poll.name != ResetPidorPoll && poll.name != ResetHeroPoll {
		return
	}
	activeChatUsers := getEnabledChatUsersByChatId(message.Chat.ID)
	var activeUsernames []string
	for _, activeChatUser := range activeChatUsers {
		activeUsernames = append(activeUsernames, FormatChatUserName(activeChatUser))
	}

	poll.applySelectedOption(update.CallbackQuery.From, ResetPollAgreedOption, ResetPollDisagreedOption)
	poll.improveVotedUserArrays(activeUsernames)
	agreedPercentage := calcPercentage(len(poll.agreedUsers), len(activeUsernames))
	disagreedPercentage := calcPercentage(len(poll.disagreedUsers), len(activeUsernames))
	poll.updateButtonsText(agreedPercentage, disagreedPercentage)

	if agreedPercentage >= config.BotResetMinPercentage {
		resetByPollName(message.Chat.ID, message.MessageID, update.CallbackQuery.ID, poll.name)
		return
	} else if disagreedPercentage >= config.BotResetMinPercentage {
		resetByPollName(message.Chat.ID, message.MessageID, update.CallbackQuery.ID, `nil`)
		return
	}
	editedMarkup := buildResetPollMarkup(poll)
	editedMsg := tgbotapi.NewEditMessageReplyMarkup(message.Chat.ID, message.MessageID, editedMarkup)
	_, _ = bot.Send(editedMsg)
}

func calcPercentage(part int, full int) int {
	if full == 0 {
		return 0
	}
	return int(float32(part) / float32(full) * 100)
}

func resetByPollName(chatId int64, msgId int, callbackId string, pollName string) {
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
	_, _ = bot.Send(deleteKeyBoardMsg)
}

func buildResetPollMarkup(resetPoll *SimplePoll) tgbotapi.InlineKeyboardMarkup {
	agreedButtonData := buildSimplePollButtonData(resetPoll, ResetPollAgreedOption)
	disagreedButtonData := buildSimplePollButtonData(resetPoll, ResetPollDisagreedOption)

	return tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(resetPoll.agreedText, agreedButtonData),
			tgbotapi.NewInlineKeyboardButtonData(resetPoll.disagreedText, disagreedButtonData),
		))
}
