package services

import (
	"github.com/LugaMuga/UOFDBot/internal/dao"
	"github.com/LugaMuga/UOFDBot/internal/locale"
	"github.com/LugaMuga/UOFDBot/internal/models"
	"github.com/LugaMuga/UOFDBot/internal/utils"
	tgbotapi "github.com/Syfaro/telegram-bot-api"
	"strconv"
	"strings"
)

const SimplePollUserDelimiter = `|`

/*
Protocol description
SIMPLE_POLL||RESET_PIDOR||yes||@michael|@mike...||@jimm|@jack...
CallbackQueryType||ResetPidorPoll||ResetPollAgreedOption||agreedUser|agreedUser...||disagreedUser|disagreedUser...
*/
type SimplePoll struct {
	name           string
	text           string
	selectedOption string
	agreedText     string
	agreedUsers    []string
	disagreedText  string
	disagreedUsers []string
}

func NewSimplePoll(pollName string) *SimplePoll {
	simplePoll := new(SimplePoll)
	simplePoll.name = pollName
	simplePoll.updateButtonsText(0, 0)
	return simplePoll
}

func (simplePoll *SimplePoll) updateButtonsText(agreedPercentage int, disagreedPercentage int) {
	simplePoll.agreedText = locale.Loc(locale.DefaultLang, `yes%`, strconv.Itoa(agreedPercentage))
	simplePoll.disagreedText = locale.Loc(locale.DefaultLang, `no%`, strconv.Itoa(disagreedPercentage))
	if len(simplePoll.agreedUsers) > 0 {
		simplePoll.agreedText += ` [` + joinUsers(simplePoll.agreedUsers) + `]`
	}
	if len(simplePoll.disagreedUsers) > 0 {
		simplePoll.disagreedText += ` [` + joinUsers(simplePoll.disagreedUsers) + `]`
	}
}

func (simplePoll *SimplePoll) applySelectedOption(user *tgbotapi.User, agreedOption string, disagreedOption string) {
	username := utils.FormatUserNameFromApi(user)
	if simplePoll.selectedOption == agreedOption {
		updateUserArrays(username, &simplePoll.agreedUsers, &simplePoll.disagreedUsers)
	} else if simplePoll.selectedOption == disagreedOption {
		updateUserArrays(username, &simplePoll.disagreedUsers, &simplePoll.agreedUsers)
	}
}

func updateUserArrays(username string, firstUsers *[]string, secondUsers *[]string) {
	if !utils.Contains(*firstUsers, username) {
		*firstUsers = append(*firstUsers, username)
	}
	if i := utils.IndexOf(*secondUsers, username); i >= 0 {
		*secondUsers = utils.Remove(*secondUsers, i)
	}
}

func (simplePoll *SimplePoll) improveVotedUserArrays(activeUsernames []string) {
	utils.RetainAll(&simplePoll.agreedUsers, activeUsernames)
	utils.RetainAll(&simplePoll.disagreedUsers, activeUsernames)
}

func joinUsers(users []string) string {
	if len(users) == 0 {
		return ``
	}
	return strings.Join(users[:], SimplePollUserDelimiter)
}

func splitUsers(optionParam string) []string {
	users := strings.Split(optionParam, SimplePollUserDelimiter)
	if len(users) == 1 && len(users[0]) == 0 {
		return []string{}
	}
	return users
}

func buildSimplePollButtonData(resetPoll *SimplePoll, option string, chatCallback *models.ChatCallback) string {
	var sb strings.Builder
	sb.WriteString(string(SimplePollType))
	sb.WriteString(CallbackQueryParamDelimiter)
	sb.WriteString(resetPoll.name)
	sb.WriteString(CallbackQueryParamDelimiter)
	sb.WriteString(option)
	sb.WriteString(CallbackQueryParamDelimiter)
	var text strings.Builder
	text.WriteString(joinUsers(resetPoll.agreedUsers))
	text.WriteString(CallbackQueryParamDelimiter)
	text.WriteString(joinUsers(resetPoll.disagreedUsers))
	chatCallback.Text = text.String()
	sb.WriteString(strconv.FormatInt(chatCallback.Id, 10))
	return sb.String()
}

func ParseSimplePollCallbackQuery(query *tgbotapi.CallbackQuery) *SimplePoll {
	simplePoll := new(SimplePoll)
	simplePoll.text = query.Message.Text
	callbackData := query.Data
	params := strings.Split(callbackData, CallbackQueryParamDelimiter)
	id, _ := strconv.Atoi(params[3])
	callback := dao.GetChatCallbackById(id)
	userParams := strings.Split(callback.Text, CallbackQueryParamDelimiter)
	simplePoll.name = params[1]
	simplePoll.selectedOption = params[2]
	simplePoll.agreedUsers = splitUsers(userParams[0])
	simplePoll.disagreedUsers = splitUsers(userParams[1])
	return simplePoll
}
