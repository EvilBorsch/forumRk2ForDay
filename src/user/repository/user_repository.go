package repository

import (
	main "../../../go"
	umodel "../models"
)

func SaveUser(user umodel.User) error {
	conn := main.GetConnection()
	query := `INSERT INTO "user" (nickname, fullname, about, email) VALUES ($1,$2,$3,$4)`
	_, err := conn.Exec(query, user.Nickname, user.Fullname, user.About, user.Email)
	return err
}
