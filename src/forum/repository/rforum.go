package repository

import (
	urepo "../../user/repository"
	"../../utills"
	fmodel "../models"
	"fmt"
)

func CreateForum(forum fmodel.Forum) (fmodel.Forum, error) {
	conn := utills.GetConnection()
	var newForum fmodel.Forum
	user, err := urepo.GetUserByNickname(forum.User_nickname)
	fmt.Println(user, err)
	if err != nil {
		return fmodel.Forum{}, err
	}
	forum.User_nickname = user.Nickname
	query := `INSERT INTO forum (title, user_nickname, slug) VALUES ($1,$2,$3) returning *`
	err = conn.Get(&newForum, query, forum.Title, forum.User_nickname, forum.Slug)

	return newForum, err
}

func GetForumBySlug(slug string) (fmodel.Forum, error) {
	conn := utills.GetConnection()
	var forum fmodel.Forum
	query := `SELECT * from forum where slug=$1`
	err := conn.Get(&forum, query, slug)
	return forum, err
}
