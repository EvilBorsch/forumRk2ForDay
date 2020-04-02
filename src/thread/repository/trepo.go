package repository

import (
	"../../utills"
	tm "../models"
)

func AddNew(newThread tm.Thread) (tm.Thread, error) {
	conn := utills.GetConnection()
	query := `INSERT INTO threads (author, created, forum, message, slug, title) VALUES ($1,$2,$3,$4,$5,$6) returning *`
	var createdThred tm.Thread
	err := conn.Get(&createdThred, query, newThread.Author, newThread.Created, newThread.Forum, newThread.Message, newThread.Slug, newThread.Title)
	return createdThred, err

}
