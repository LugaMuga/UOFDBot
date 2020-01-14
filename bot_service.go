package main

import (
	tgbotapi "github.com/Syfaro/telegram-bot-api"
	"math/rand"
	"time"
)

const resetPidorApproval string = `RESET_PIDOR_SCORES`
const resetHeroApproval string = `RESET_HERO_SCORES`
const resetCancellation string = `RESET_CANCEL`

func play(chatUsers []ChatUser) int {
	if len(chatUsers) <= 0 {
		return -1
	} else {
		rand.Seed(time.Now().UnixNano())
		return rand.Intn(len(chatUsers))
	}
}

func register(message tgbotapi.Message) {
	chatUser := findChatUserByUserIdAndChatId(message.From.ID, message.Chat.ID)
	username := FormatUserName(message.From.UserName, message.From.FirstName, message.From.LastName)
	if chatUser != nil && chatUser.enabled {
		SendMessage(message.Chat.ID, loc(defaultLang, `user_already_registered`, username))
		return
	}
	if chatUser == nil {
		chatUser = new(ChatUser)
	}
	chatUser.fillFromMessage(message)
	chatUser.enabled = true
	SaveOrUpdateChatUser(*chatUser)
	SendMessage(message.Chat.ID, loc(defaultLang, `user_registered`, username))
}

func delete(message tgbotapi.Message) {
	chatUser := findChatUserByUserIdAndChatId(message.From.ID, message.Chat.ID)
	username := FormatUserName(message.From.UserName, message.From.FirstName, message.From.LastName)
	if chatUser == nil || !chatUser.enabled {
		SendMessage(message.Chat.ID, loc(defaultLang, `user_not_participating`, username))
		return
	}
	chatUser.fillFromMessage(message)
	chatUser.enabled = false
	UpdateChatUserStatus(*chatUser)
	SendMessage(message.Chat.ID, loc(defaultLang, `user_deleted`, username))
}

func pidor(chatId int64) {
	activePidor := FindActivePidorByChatId(chatId)
	if activePidor != nil {
		msg := FormatActivePidorWinner(*activePidor)
		SendMessage(chatId, msg)
		return
	}
	chatUsers := getEnabledChatUsersByChatId(chatId)
	winnerIndex := play(chatUsers)
	if winnerIndex < 0 {
		SendMessage(chatId, loc(defaultLang, `at_least_one_user`))
		return
	}
	chatUsers[winnerIndex].pidorScore += 1
	chatUsers[winnerIndex].pidorLastTimestamp = nowUnix()
	UpdateChatUserPidorWins(chatUsers[winnerIndex])
	msg := FormatPidorWinner(chatUsers[winnerIndex])
	SendMessage(chatId, msg)
}

func pidorList(chatId int64) {
	chatUsers := GetPidorListScoresByChatId(chatId)
	msg := FormatListOfPidors(chatUsers)
	SendMessage(chatId, msg)
}

func resetPidor(chatId int64) {
	ResetPidorScoreByChatId(chatId)
	gameName := loc(defaultLang, `pidor_of_day`)
	msg := loc(defaultLang, `stat_reset`, gameName)
	SendMessage(chatId, msg)
}

func hero(chatId int64) {
	activeHero := FindActiveHeroByChatId(chatId)
	if activeHero != nil {
		msg := FormatActiveHeroWinner(*activeHero)
		SendMessage(chatId, msg)
		return
	}
	chatUsers := getEnabledChatUsersByChatId(chatId)
	winnerIndex := play(chatUsers)
	if winnerIndex < 0 {
		SendMessage(chatId, loc(defaultLang, `at_least_one_user`))
		return
	}
	chatUsers[winnerIndex].heroScore += 1
	chatUsers[winnerIndex].heroLastTimestamp = nowUnix()
	UpdateChatUserHeroWins(chatUsers[winnerIndex])
	msg := FormatHeroWinner(chatUsers[winnerIndex])
	SendMessage(chatId, msg)
}

func heroList(chatId int64) {
	chatUsers := GetHeroListScoresByChatId(chatId)
	msg := FormatListOfHeros(chatUsers)
	SendMessage(chatId, msg)
}

func resetHero(chatId int64) {
	ResetHeroScoreByChatId(chatId)
	gameName := loc(defaultLang, `hero_of_day`)
	msg := loc(defaultLang, `stat_reset`, gameName)
	SendMessage(chatId, msg)
}

func resetApproval(chatId int64, approvalOption string) {
	msg := tgbotapi.NewMessage(chatId, loc(defaultLang, `a_u_sure`))
	msg.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(loc(defaultLang, `yes`), approvalOption),
			tgbotapi.NewInlineKeyboardButtonData(loc(defaultLang, `no`), resetCancellation),
		))
	_, _ = bot.Send(msg)
}

func run(chatId int64) {
	pidor(chatId)
	hero(chatId)
}

func list(chatId int64) {
	pidorList(chatId)
	heroList(chatId)
}
