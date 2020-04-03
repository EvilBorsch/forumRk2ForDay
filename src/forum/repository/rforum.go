package repository

import (
	tmodel "../../thread/models"
	urepo "../../user/repository"
	"../../utills"
	fmodel "../models"
	"fmt"
	"time"
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

func GetThreadsByForumSlug(forumSlug string, isDesc string, limit string, since string) ([]tmodel.Thread, error, bool) {
	conn := utills.GetConnection()
	var thread []tmodel.Thread
	query := `SELECT * from threads where forum=$1 ORDER BY created`
	_, err := GetForumBySlug(forumSlug)
	if err != nil {
		return []tmodel.Thread{}, nil, false
	}
	if since != "" {
		fmt.Println(since)
		layout := "2006-01-02T15:04:05.000Z"
		t, _ := time.Parse(layout, since)
		t = t.UTC().Local()
		tString := "'" + t.String()[:23] + "'"
		fmt.Println(tString)
		query = `SELECT * from threads where forum=$1 and created<=` + tString + ` ORDER BY created`
		if isDesc == "false" || isDesc == "" {
			query = `SELECT * from threads where forum=$1 and created>=` + tString + ` ORDER BY created`
		}
	}
	if isDesc == "true" {
		query += " DESC"
	}
	if limit != "" {
		query = query + " LIMIT " + limit
	}
	err = conn.Select(&thread, query, forumSlug)
	return thread, err, true
}

func IncrementFieldBySlug(fieldName string, slug string) error {
	query := fmt.Sprintf(`UPDATE forum SET %s =%s + 1 WHERE slug=$1`, fieldName, fieldName)
	fmt.Println(query)
	conn := utills.GetConnection()
	_, err := conn.Exec(query, slug)
	return err
}
