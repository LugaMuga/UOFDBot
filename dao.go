package main

import "log"

func SaveOrUpdateChatUser(chatUser ChatUser) {
	if chatUser.Id == 0 {
		InsertChatUser(chatUser)
		return
	}
	UpdateChatUserStatus(chatUser)
}

func SaveOrUpdateChatCallback(chatCallback ChatCallback) int64 {
	if chatCallback.Id == 0 {
		return InsertChatCallback(chatCallback)
	}
	UpdateChatCallbackText(chatCallback)
	return chatCallback.Id
}

func InsertChatUser(chatUser ChatUser) {
	tx, err := DB.Begin()
	if err != nil {
		log.Fatal(err)
	}
	stmt, err := tx.Prepare(`INSERT INTO chat_user(
                 chat_id,
                 user_id,
                 username,
                 user_first_name,
                 user_last_name,
                 enabled)
                 values(?, ?, ?, ?, ?, ?)`)
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()
	_, err = stmt.Exec(
		chatUser.ChatId,
		chatUser.UserId,
		chatUser.Username,
		chatUser.UserFirstName,
		chatUser.UserLastName,
		chatUser.Enabled)
	if err != nil {
		log.Fatal(err)
	}
	tx.Commit()
}

func InsertChatCallback(chatCallback ChatCallback) int64 {
	tx, err := DB.Begin()
	if err != nil {
		log.Fatal(err)
	}
	stmt, err := tx.Prepare(`INSERT INTO chat_callback(
                 chat_id,
                 text,
                 create_timestamp)
                 values(?, ?, ?)`)
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()

	res, err := stmt.Exec(
		chatCallback.ChatId,
		chatCallback.Text,
		chatCallback.CreateTimestamp)
	if err != nil {
		log.Fatal(err)
	}
	tx.Commit()
	id, err := res.LastInsertId()
	return id
}

func UpdateChatCallbackText(chatCallback ChatCallback) {
	tx, err := DB.Begin()
	if err != nil {
		log.Fatal(err)
	}
	stmt, err := tx.Prepare(`UPDATE chat_callback SET text = ? WHERE chat_id = ? AND id = ?`)
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()
	_, err = stmt.Exec(
		chatCallback.Text,
		chatCallback.ChatId,
		chatCallback.Id)
	if err != nil {
		log.Fatal(err)
	}
	err = tx.Commit()
	if err != nil {
		log.Fatal(err)
	}
}

func DeleteChatCallback(id int64, chatId int64) {
	tx, err := DB.Begin()
	if err != nil {
		log.Fatal(err)
	}
	stmt, err := tx.Prepare(`DELETE FROM chat_callback WHERE id = ? and chat_id = ?`)
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()
	_, err = stmt.Exec(
		id,
		chatId)
	if err != nil {
		log.Fatal(err)
	}
	err = tx.Commit()
	if err != nil {
		log.Fatal(err)
	}
}

func UpdateChatUserStatus(chatUser ChatUser) {
	tx, err := DB.Begin()
	if err != nil {
		log.Fatal(err)
	}
	stmt, err := tx.Prepare(`UPDATE chat_user SET enabled = ? WHERE chat_id = ? AND user_id = ?`)
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()
	_, err = stmt.Exec(
		chatUser.Enabled,
		chatUser.ChatId,
		chatUser.UserId)
	if err != nil {
		log.Fatal(err)
	}
	err = tx.Commit()
	if err != nil {
		log.Fatal(err)
	}
}

func UpdateChatUserPidorWins(chatUser ChatUser) {
	updateChatUserScore(chatUser.Id, chatUser.PidorScore, `pidor_score`, chatUser.PidorLastTimestamp, `pidor_last_timestamp`)
}

func UpdateChatUserHeroWins(chatUser ChatUser) {
	updateChatUserScore(chatUser.Id, chatUser.HeroScore, `hero_score`, chatUser.HeroLastTimestamp, `hero_last_timestamp`)
}

func updateChatUserScore(chatUserId int, score int, scoreField string, timestamp int64, timestampField string) {
	tx, err := DB.Begin()
	if err != nil {
		log.Fatal(err)
	}
	stmt, err := tx.Prepare(`UPDATE chat_user SET ` + scoreField + ` = ?, ` + timestampField + `= ? WHERE id = ? `)
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()
	_, err = stmt.Exec(score, timestamp, chatUserId)
	if err != nil {
		log.Fatal(err)
	}
	err = tx.Commit()
	if err != nil {
		log.Fatal(err)
	}
}

func findChatUserByUserIdAndChatId(userId int, chatId int64) *ChatUser {
	stmt, err := DB.Prepare(`
		SELECT 	id,
		       	chat_id,
				user_id,
				username,
				user_first_name,
				user_last_name,
				enabled,
				pidor_score,
				pidor_last_timestamp,
				hero_score,
				hero_last_timestamp
		FROM chat_user
		WHERE user_id = ? AND chat_id = ?`)
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()
	rows, err := stmt.Query(userId, chatId)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	if !rows.Next() {
		return nil
	}
	chatUser := new(ChatUser)
	err = rows.Scan(
		&chatUser.Id,
		&chatUser.ChatId,
		&chatUser.UserId,
		&chatUser.Username,
		&chatUser.UserFirstName,
		&chatUser.UserLastName,
		&chatUser.Enabled,
		&chatUser.PidorScore,
		&chatUser.PidorLastTimestamp,
		&chatUser.HeroScore,
		&chatUser.HeroLastTimestamp)
	if err != nil {
		log.Fatal(err)
	}
	return chatUser
}

func FindActivePidorByChatId(chatId int64) *ChatUser {
	return findEnabledChatUserWonInGameIntervalByChatId(chatId, `pidor_last_timestamp`)
}

func FindActiveHeroByChatId(chatId int64) *ChatUser {
	return findEnabledChatUserWonInGameIntervalByChatId(chatId, `hero_last_timestamp`)
}

func findEnabledChatUserWonInGameIntervalByChatId(chatId int64, timestampField string) *ChatUser {
	stmt, err := DB.Prepare(`
		SELECT user_id,
			   username,
			   user_first_name,
			   user_last_name,
			   pidor_score,
			   hero_score
		FROM chat_user
		WHERE chat_id = ? AND enabled AND ` + timestampField + `> ?`)
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()
	rows, err := stmt.Query(chatId, getLastMidnight())
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	if !rows.Next() {
		return nil
	}
	chatUser := new(ChatUser)
	err = rows.Scan(
		&chatUser.UserId,
		&chatUser.Username,
		&chatUser.UserFirstName,
		&chatUser.UserLastName,
		&chatUser.PidorScore,
		&chatUser.HeroScore)
	if err != nil {
		log.Fatal(err)
	}
	return chatUser
}

func getEnabledChatUsersByChatId(chatId int64) []ChatUser {
	stmt, err := DB.Prepare(`
		SELECT id,
			   user_id,
			   username,
			   user_first_name,
			   user_last_name,
			   pidor_score,
			   hero_score
		FROM chat_user
		WHERE chat_id = ? AND enabled
		ORDER BY user_id`)
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()
	rows, err := stmt.Query(chatId)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	var chatUsers []ChatUser
	for rows.Next() {
		chatUser := new(ChatUser)
		err = rows.Scan(
			&chatUser.Id,
			&chatUser.UserId,
			&chatUser.Username,
			&chatUser.UserFirstName,
			&chatUser.UserLastName,
			&chatUser.PidorScore,
			&chatUser.HeroScore)
		if err != nil {
			log.Fatal(err)
		}
		chatUsers = append(chatUsers, *chatUser)
	}
	return chatUsers
}

func getChatCallbackById(id int) *ChatCallback {
	stmt, err := DB.Prepare(`
		SELECT id,
			   chat_id,
			   text,
			   create_timestamp
		FROM chat_callback
		WHERE id = ?`)
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()
	rows, err := stmt.Query(id)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	if !rows.Next() {
		return nil
	}
	chatCallback := new(ChatCallback)
	err = rows.Scan(
		&chatCallback.Id,
		&chatCallback.ChatId,
		&chatCallback.Text,
		&chatCallback.CreateTimestamp)
	if err != nil {
		log.Fatal(err)
	}
	return chatCallback
}

func GetPidorListScoresByChatId(chatId int64) []ChatUser {
	return getEnabledScoreListByChatId(chatId, `pidor_score`, `pidor_last_timestamp`)
}

func GetHeroListScoresByChatId(chatId int64) []ChatUser {
	return getEnabledScoreListByChatId(chatId, `hero_score`, `hero_last_timestamp`)
}

func getEnabledScoreListByChatId(chatId int64, scoreField string, timestampField string) []ChatUser {
	stmt, err := DB.Prepare(`
		SELECT user_id,
			   username,
			   user_first_name,
			   user_last_name,
			   pidor_score,
		       pidor_last_timestamp,
			   hero_score,
		       hero_last_timestamp
		FROM chat_user
		WHERE chat_id = ? AND enabled
		ORDER BY ` + scoreField + ` DESC, ` + timestampField + ` DESC`)
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()
	rows, err := stmt.Query(chatId)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	var chatUsers []ChatUser
	for rows.Next() {
		chatUser := new(ChatUser)
		err = rows.Scan(
			&chatUser.UserId,
			&chatUser.Username,
			&chatUser.UserFirstName,
			&chatUser.UserLastName,
			&chatUser.PidorScore,
			&chatUser.PidorLastTimestamp,
			&chatUser.HeroScore,
			&chatUser.HeroLastTimestamp)
		if err != nil {
			log.Fatal(err)
		}
		chatUsers = append(chatUsers, *chatUser)
	}
	return chatUsers
}

func ResetPidorScoreByChatId(chatId int64) {
	resetScoreByChatId(chatId, `pidor_score`, `pidor_last_timestamp`)
}

func ResetHeroScoreByChatId(chatId int64) {
	resetScoreByChatId(chatId, `hero_score`, `hero_last_timestamp`)
}

func resetScoreByChatId(chatId int64, scoreField string, timestampField string) {
	tx, err := DB.Begin()
	if err != nil {
		log.Fatal(err)
	}
	stmt, err := tx.Prepare(`UPDATE chat_user SET ` + scoreField + ` = 0, ` + timestampField + `=0  WHERE chat_id = ?`)
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()
	_, err = stmt.Exec(chatId)
	if err != nil {
		log.Fatal(err)
	}
	err = tx.Commit()
	if err != nil {
		log.Fatal(err)
	}
}
