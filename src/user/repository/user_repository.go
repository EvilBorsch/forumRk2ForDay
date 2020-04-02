package repository

import (
	"../../utills"
	umodel "../models"
)

func SaveUser(user umodel.User) error {
	conn := utills.GetConnection()
	query := `INSERT INTO "user" (nickname, fullname, about, email) VALUES ($1,$2,$3,$4)`
	_, err := conn.Exec(query, user.Nickname, user.Fullname, user.About, user.Email)
	return err
}

func GetUserByNickname(nickname string) (umodel.User, error) {
	conn := utills.GetConnection()
	query := `SELECT * from "user" WHERE nickname=$1`
	var user umodel.User
	err := conn.Get(&user, query, nickname)
	return user, err
}
