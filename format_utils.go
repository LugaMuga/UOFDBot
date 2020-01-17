package main

import (
	tgbotapi "github.com/Syfaro/telegram-bot-api"
	"strconv"
	"strings"
	"time"
)

type winFunc func(int) int
type comparatorFunc func(ChatUser, int64) int64

func FormatUserNameFromApi(user *tgbotapi.User) string {
	return FormatUserName(user.UserName, user.FirstName, user.LastName)
}

func FormatChatUserName(chatUser ChatUser) string {
	return FormatUserName(chatUser.Username, chatUser.UserFirstName, chatUser.UserLastName)
}

func FormatUserName(username string, firstName string, lastName string) string {
	if len(username) > 0 {
		return `@` + username
	}
	return firstName + ` ` + lastName
}

func FormatActivePidorWinner(chatUser ChatUser) string {
	return formatWinnerMsg(chatUser, "\U0001F308"+loc(defaultLang, `today_pidor_selected`)+" - ")
}

func FormatActiveHeroWinner(chatUser ChatUser) string {
	return formatWinnerMsg(chatUser, "\U0001F31F"+loc(defaultLang, `today_hero_selected`)+" - ")
}

func FormatPidorWinner(chatUser ChatUser) string {
	return formatWinnerMsg(chatUser, "\U0001F308"+loc(defaultLang, `today_pidor`)+" - ")
}

func FormatHeroWinner(chatUser ChatUser) string {
	return formatWinnerMsg(chatUser, "\U0001F31F"+loc(defaultLang, `today_hero`)+" - ")
}

func formatWinnerMsg(chatUser ChatUser, title string) string {
	var sb strings.Builder
	sb.WriteString(title)
	sb.WriteString(chatUser.UserFirstName + ` ` + chatUser.UserLastName)
	if len(chatUser.Username) > 0 {
		sb.WriteString(` (@` + chatUser.Username + `)`)
	}
	return sb.String()
}

func FormatListOfPidors(chatUsers []ChatUser) string {
	lastRunPidorComparator := func(chatUser ChatUser, lastRun int64) int64 {
		if chatUser.PidorLastTimestamp > 0 && chatUser.PidorLastTimestamp > lastRun {
			return chatUser.PidorLastTimestamp
		}
		return lastRun
	}
	lastTimeRun := findLastTimeRun(chatUsers, lastRunPidorComparator)

	getNumberOfWins := func(i int) int {
		return chatUsers[i].PidorScore
	}
	resultsMsg := loc(defaultLang, `results_by_game`, loc(defaultLang, `pidor_of_day`))
	return formatListOfGames(chatUsers, resultsMsg+" \U0001F308 "+lastTimeRun, getNumberOfWins)
}

func FormatListOfHeros(chatUsers []ChatUser) string {
	lastRunHeroComparator := func(chatUser ChatUser, lastRun int64) int64 {
		if chatUser.HeroLastTimestamp > 0 && chatUser.HeroLastTimestamp > lastRun {
			return chatUser.HeroLastTimestamp
		}
		return lastRun
	}
	lastTimeRun := findLastTimeRun(chatUsers, lastRunHeroComparator)

	getNumberOfWins := func(i int) int {
		return chatUsers[i].HeroScore
	}
	resultsMsg := loc(defaultLang, `results_by_game`, loc(defaultLang, `hero_of_day`))
	return formatListOfGames(chatUsers, resultsMsg+" \U0001F31F "+lastTimeRun, getNumberOfWins)
}

func formatListOfGames(chatUsers []ChatUser, title string, getNumberOfWins winFunc) string {
	var sb strings.Builder
	sb.WriteString(title + "\n")
	for i := 0; i < len(chatUsers); i++ {
		sb.WriteString(strconv.Itoa(i+1) + `) `)
		sb.WriteString(chatUsers[i].UserFirstName + ` ` + chatUsers[i].UserLastName)
		if len(chatUsers[i].Username) > 0 {
			sb.WriteString(` (@` + chatUsers[i].Username + `)`)
		}
		timesText := locPlural(defaultLang, "win_times", getNumberOfWins(i), getNumberOfWins(i))
		sb.WriteString(` - ` + timesText)
		sb.WriteString("\n")
	}
	return sb.String()
}

func findLastTimeRun(chatUsers []ChatUser, comparator comparatorFunc) string {
	if len(config.BotTimeLayout) == 0 {
		return ``
	}
	var lastRun int64 = 0
	for _, chatUser := range chatUsers {
		lastRun = comparator(chatUser, lastRun)
	}
	if lastRun == 0 {
		return loc(defaultLang, `last_run`, loc(defaultLang, `never`))
	}
	lastRunTime := time.Unix(lastRun, 0)
	return loc(defaultLang, `last_run`, lastRunTime.Format(config.BotTimeLayout))
}
