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

func GetUserByNicknameAndEmail(nickname string, email string) ([]umodel.User, error) {
	conn := utills.GetConnection()
	query := `SELECT * from "user" WHERE nickname=$1 OR email=$2`
	var userlist []umodel.User
	err := conn.Select(&userlist, query, nickname, email)
	return userlist, err
}

func UpdateUser(oldNickname string, newUser umodel.User) (umodel.User, error) {
	conn := utills.GetConnection()
	query := `UPDATE "user" SET fullname=$1,about=$2,email=$3 WHERE nickname=$4 RETURNING *`
	var updatedUser umodel.User
	err := conn.Get(&updatedUser, query, newUser.Fullname, newUser.About, newUser.Email, oldNickname)
	return updatedUser, err

}
