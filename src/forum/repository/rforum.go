package repository

import (
	"../../utills"
	fmodel "../models"
)

func CreateForum(forum fmodel.Forum) (fmodel.Forum, error) {
	conn := utills.GetConnection()
	var newForum fmodel.Forum
	query := `INSERT INTO forum (title, user_nickname, slug) VALUES ($1,$2,$3) returning *`
	err := conn.Get(&newForum, query, forum.Title, forum.User_nickname, forum.Slug)
	return newForum, err
}
