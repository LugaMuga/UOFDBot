package main

import (
	"strconv"
	"strings"
)

type winFunc func(int) string

func FormatUserName(username string, firstName string, lastName string) string {
	if len(username) > 0 {
		return `@` + username
	}
	return firstName + ` ` + lastName
}

func FormatActivePidorWinner(chatUser ChatUser) string {
	return formatWinnerMsg(chatUser, "\U0001F308Сегодня ПИДОР дня уже был выбран - ")
}

func FormatActiveHeroWinner(chatUser ChatUser) string {
	return formatWinnerMsg(chatUser, "\U0001F31FСегодня ГЕРОЙ дня уже был выбран - ")
}

func FormatPidorWinner(chatUser ChatUser) string {
	return formatWinnerMsg(chatUser, "\U0001F308Сегодня ПИДОР дня - ")
}

func FormatHeroWinner(chatUser ChatUser) string {
	return formatWinnerMsg(chatUser, "\U0001F31FСегодня ГЕРОЙ дня - ")
}

func formatWinnerMsg(chatUser ChatUser, title string) string {
	var sb strings.Builder
	sb.WriteString(title)
	sb.WriteString(chatUser.userFirstName + ` ` + chatUser.userLastName)
	if len(chatUser.username) > 0 {
		sb.WriteString(` (@` + chatUser.username + `)`)
	}
	return sb.String()
}

func FormatListOfPidors(chatUsers []ChatUser) string {
	getNumberOfWins := func(i int) string {
		return strconv.Itoa(chatUsers[i].pidorScore)
	}
	return formatListOfGames(chatUsers, "Итоги 'пидора дня' \U0001F308", getNumberOfWins)
}

func FormatListOfHeros(chatUsers []ChatUser) string {
	getNumberOfWins := func(i int) string {
		return strconv.Itoa(chatUsers[i].heroScore)
	}
	return formatListOfGames(chatUsers, "Итоги 'героя дня' \U0001F31F", getNumberOfWins)
}

func formatListOfGames(chatUsers []ChatUser, title string, getNumberOfWins winFunc) string {
	var sb strings.Builder
	sb.WriteString(title + "\n")
	for i := 0; i < len(chatUsers); i++ {
		sb.WriteString(strconv.Itoa(i+1) + `) `)
		sb.WriteString(chatUsers[i].userFirstName + ` ` + chatUsers[i].userLastName)
		if len(chatUsers[i].username) > 0 {
			sb.WriteString(` (@` + chatUsers[i].username + `)`)
		}
		sb.WriteString(` - ` + getNumberOfWins(i) + ` раз(а)`)
		sb.WriteString("\n")
	}
	return sb.String()
}
