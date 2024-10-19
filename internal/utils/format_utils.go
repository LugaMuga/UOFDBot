package utils

import (
	"github.com/LugaMuga/UOFDBot/internal/config"
	"github.com/LugaMuga/UOFDBot/internal/locale"
	"github.com/LugaMuga/UOFDBot/internal/models"
	tgbotapi "github.com/Syfaro/telegram-bot-api"
	"strconv"
	"strings"
	"time"
)

type winFunc func(int) int
type comparatorFunc func(models.ChatUser, int64) int64

func FormatUserNameFromApi(user *tgbotapi.User) string {
	return FormatUserName(user.UserName, user.FirstName, user.LastName)
}

func FormatChatUserName(chatUser models.ChatUser) string {
	return FormatUserName(chatUser.Username, chatUser.UserFirstName, chatUser.UserLastName)
}

func FormatUserName(username string, firstName string, lastName string) string {
	if len(username) > 0 {
		return `@` + username
	}
	return firstName + ` ` + lastName
}

func FormatActivePidorWinner(chatUser models.ChatUser) string {
	return formatWinnerMsg(chatUser, "\U0001F308"+locale.Loc(locale.DefaultLang, `today_pidor_selected`)+" - ")
}

func FormatActiveHeroWinner(chatUser models.ChatUser) string {
	return formatWinnerMsg(chatUser, "\U0001F31F"+locale.Loc(locale.DefaultLang, `today_hero_selected`)+" - ")
}

func FormatPidorWinner(chatUser models.ChatUser) string {
	return formatWinnerMsg(chatUser, "\U0001F308"+locale.Loc(locale.DefaultLang, `today_pidor`)+" - ")
}

func FormatHeroWinner(chatUser models.ChatUser) string {
	return formatWinnerMsg(chatUser, "\U0001F31F"+locale.Loc(locale.DefaultLang, `today_hero`)+" - ")
}

func formatWinnerMsg(chatUser models.ChatUser, title string) string {
	var sb strings.Builder
	sb.WriteString(title)
	sb.WriteString(chatUser.UserFirstName + ` ` + chatUser.UserLastName)
	if len(chatUser.Username) > 0 {
		sb.WriteString(` (@` + chatUser.Username + `)`)
	}
	return sb.String()
}

func FormatListOfPidors(chatUsers []models.ChatUser) string {
	lastRunPidorComparator := func(chatUser models.ChatUser, lastRun int64) int64 {
		if chatUser.PidorLastTimestamp > 0 && chatUser.PidorLastTimestamp > lastRun {
			return chatUser.PidorLastTimestamp
		}
		return lastRun
	}
	lastTimeRun := findLastTimeRun(chatUsers, lastRunPidorComparator)

	getNumberOfWins := func(i int) int {
		return chatUsers[i].PidorScore
	}
	resultsMsg := locale.Loc(locale.DefaultLang, `results_by_game`, locale.Loc(locale.DefaultLang, `pidor_of_day`))
	return formatListOfGames(chatUsers, resultsMsg+" \U0001F308 "+lastTimeRun, getNumberOfWins)
}

func FormatListOfHeros(chatUsers []models.ChatUser) string {
	lastRunHeroComparator := func(chatUser models.ChatUser, lastRun int64) int64 {
		if chatUser.HeroLastTimestamp > 0 && chatUser.HeroLastTimestamp > lastRun {
			return chatUser.HeroLastTimestamp
		}
		return lastRun
	}
	lastTimeRun := findLastTimeRun(chatUsers, lastRunHeroComparator)

	getNumberOfWins := func(i int) int {
		return chatUsers[i].HeroScore
	}
	resultsMsg := locale.Loc(locale.DefaultLang, `results_by_game`, locale.Loc(locale.DefaultLang, `hero_of_day`))
	return formatListOfGames(chatUsers, resultsMsg+" \U0001F31F "+lastTimeRun, getNumberOfWins)
}

func formatListOfGames(chatUsers []models.ChatUser, title string, getNumberOfWins winFunc) string {
	var sb strings.Builder
	sb.WriteString(title + "\n")
	for i := 0; i < len(chatUsers); i++ {
		sb.WriteString(strconv.Itoa(i+1) + `) `)
		sb.WriteString(chatUsers[i].UserFirstName + ` ` + chatUsers[i].UserLastName)
		if len(chatUsers[i].Username) > 0 {
			sb.WriteString(` (@` + chatUsers[i].Username + `)`)
		}
		timesText := locale.LocPlural(locale.DefaultLang, "win_times", getNumberOfWins(i), getNumberOfWins(i))
		sb.WriteString(` - ` + timesText)
		sb.WriteString("\n")
	}
	return sb.String()
}

func findLastTimeRun(chatUsers []models.ChatUser, comparator comparatorFunc) string {
	if len(config.Config.BotTimeLayout) == 0 {
		return ``
	}
	var lastRun int64 = 0
	for _, chatUser := range chatUsers {
		lastRun = comparator(chatUser, lastRun)
	}
	if lastRun == 0 {
		return locale.Loc(locale.DefaultLang, `last_run`, locale.Loc(locale.DefaultLang, `never`))
	}
	lastRunTime := time.Unix(lastRun, 0)
	return locale.Loc(locale.DefaultLang, `last_run`, lastRunTime.Format(config.Config.BotTimeLayout))
}
