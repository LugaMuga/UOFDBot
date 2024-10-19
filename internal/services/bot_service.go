package services

import (
	"crypto/rand"
	"github.com/LugaMuga/UOFDBot/internal/bot"
	"github.com/LugaMuga/UOFDBot/internal/dao"
	"github.com/LugaMuga/UOFDBot/internal/locale"
	"github.com/LugaMuga/UOFDBot/internal/models"
	"github.com/LugaMuga/UOFDBot/internal/utils"
	tgbotapi "github.com/Syfaro/telegram-bot-api"
	"log"
	"math/big"
)

type CallbackQueryType string

const SimplePollType CallbackQueryType = `SIMPLE_POLL`
const CallbackQueryParamDelimiter = `||`

func play(chatUsers []models.ChatUser) int64 {
	if len(chatUsers) <= 0 {
		return -1
	} else {
		numberOfChatUsers := int64(len(chatUsers))
		r, _ := rand.Int(rand.Reader, big.NewInt(numberOfChatUsers))
		return r.Int64()
	}
}

func Register(message tgbotapi.Message) {
	chatUser := dao.FindChatUserByUserIdAndChatId(message.From.ID, message.Chat.ID)
	username := utils.FormatUserName(message.From.UserName, message.From.FirstName, message.From.LastName)
	if chatUser != nil && chatUser.Enabled {
		bot.SendMessage(message.Chat.ID, locale.Loc(locale.DefaultLang, `user_already_registered`, username))
		return
	}
	if chatUser == nil {
		chatUser = new(models.ChatUser)
	}
	chatUser.Fill(message.Chat.ID, message.From)
	chatUser.Enabled = true
	dao.SaveOrUpdateChatUser(*chatUser)
	bot.SendMessage(message.Chat.ID, locale.Loc(locale.DefaultLang, `user_registered`, username))
}

func Delete(chatId int64, user *tgbotapi.User) {
	chatUser := dao.FindChatUserByUserIdAndChatId(user.ID, chatId)
	username := utils.FormatUserName(user.UserName, user.FirstName, user.LastName)
	if chatUser == nil || !chatUser.Enabled {
		bot.SendMessage(chatId, locale.Loc(locale.DefaultLang, `user_not_participating`, username))
		return
	}
	chatUser.Fill(chatId, user)
	chatUser.Enabled = false
	dao.UpdateChatUserStatus(*chatUser)
	bot.SendMessage(chatId, locale.Loc(locale.DefaultLang, `user_deleted`, username))
}

func Pidor(chatId int64) {
	activePidor := dao.FindActivePidorByChatId(chatId)
	if activePidor != nil {
		msg := utils.FormatActivePidorWinner(*activePidor)
		bot.SendMessage(chatId, msg)
		return
	}
	chatUsers := dao.GetEnabledChatUsersByChatId(chatId)
	winnerIndex := play(chatUsers)
	if winnerIndex < 0 {
		bot.SendMessage(chatId, locale.Loc(locale.DefaultLang, `at_least_one_user`))
		return
	}
	chatUsers[winnerIndex].PidorScore += 1
	chatUsers[winnerIndex].PidorLastTimestamp = utils.NowUnix()
	dao.UpdateChatUserPidorWins(chatUsers[winnerIndex])
	msg := utils.FormatPidorWinner(chatUsers[winnerIndex])
	bot.SendMessage(chatId, msg)
}

func UpdateUsers(chatId int64) {
	chatUsers := dao.GetEnabledChatUsersByChatId(chatId)
	for _, value := range chatUsers {
		chatConfig := tgbotapi.ChatConfigWithUser{
			ChatID: chatId,
			UserID: value.UserId,
		}
		updateUser(chatId, chatConfig, value)
	}
	msg := locale.Loc(locale.DefaultLang, `update_users`)
	bot.SendMessage(chatId, msg)
}

func updateUser(chatId int64, chatConfig tgbotapi.ChatConfigWithUser, user models.ChatUser) {
	userTemp, err := bot.Bot.GetChatMember(chatConfig)
	if err != nil {
		log.Printf("user not found userId: %d, chatId: %d, err: %q", chatConfig.UserID, chatId, err)
		return
	}
	if user.Username != userTemp.User.UserName {
		user.Username = userTemp.User.UserName
		user.ChatId = chatId
		dao.UpdateChatUserUsername(user)
	}
}

func PidorList(chatId int64) {
	chatUsers := dao.GetPidorListScoresByChatId(chatId)
	msg := utils.FormatListOfPidors(chatUsers)
	bot.SendMessage(chatId, msg)
}

func resetPidor(chatId int64) {
	dao.ResetPidorScoreByChatId(chatId)
	gameName := locale.Loc(locale.DefaultLang, `pidor_of_day`)
	msg := locale.Loc(locale.DefaultLang, `stat_reset`, gameName)
	bot.SendMessage(chatId, msg)
}

func Hero(chatId int64) {
	activeHero := dao.FindActiveHeroByChatId(chatId)
	if activeHero != nil {
		msg := utils.FormatActiveHeroWinner(*activeHero)
		bot.SendMessage(chatId, msg)
		return
	}
	chatUsers := dao.GetEnabledChatUsersByChatId(chatId)
	winnerIndex := play(chatUsers)
	if winnerIndex < 0 {
		bot.SendMessage(chatId, locale.Loc(locale.DefaultLang, `at_least_one_user`))
		return
	}
	chatUsers[winnerIndex].HeroScore += 1
	chatUsers[winnerIndex].HeroLastTimestamp = utils.NowUnix()
	dao.UpdateChatUserHeroWins(chatUsers[winnerIndex])
	msg := utils.FormatHeroWinner(chatUsers[winnerIndex])
	bot.SendMessage(chatId, msg)
}

func HeroList(chatId int64) {
	chatUsers := dao.GetHeroListScoresByChatId(chatId)
	msg := utils.FormatListOfHeros(chatUsers)
	bot.SendMessage(chatId, msg)
}

func resetHero(chatId int64) {
	dao.ResetHeroScoreByChatId(chatId)
	gameName := locale.Loc(locale.DefaultLang, `hero_of_day`)
	msg := locale.Loc(locale.DefaultLang, `stat_reset`, gameName)
	bot.SendMessage(chatId, msg)
}

func Run(chatId int64) {
	Pidor(chatId)
	Hero(chatId)
}

func List(chatId int64) {
	PidorList(chatId)
	HeroList(chatId)
}
