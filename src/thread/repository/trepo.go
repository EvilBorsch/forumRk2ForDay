package repository

import (
	frepo "../../forum/repository"
	"../../utills"
	tm "../models"
	"fmt"
	"github.com/jmoiron/sqlx"
)

func AddNew(newThread tm.Thread) (tm.Thread, error) {
	conn := utills.GetConnection()
	query := `INSERT INTO threads (author, created, forum, message, title,slug) VALUES ($1,$2,$3,$4,$5,nullif($6,'')) returning *`
	var createdThred tm.Thread
	err := conn.Get(&createdThred, query, newThread.Author, newThread.Created, newThread.Forum, newThread.Message, newThread.Title, newThread.Slug)
	createdThred.Title = newThread.Title
	if err == nil {
		err2 := frepo.IncrementFieldBySlug("threads", newThread.Forum)
		fmt.Println("new err: ", err2)
	}
	return createdThred, err

}

func GetThreadBySlug(tx *sqlx.Tx, threadSlug string) (tm.Thread, error) {
	var thread tm.Thread
	query := `SELECT * FROM threads where slug=$1`
	err := tx.Get(&thread, query, threadSlug)
	return thread, err
}

func GetThreadByID(tx *sqlx.Tx, id int) (tm.Thread, error) {
	var thread tm.Thread
	query := `SELECT * FROM threads where id=$1`
	err := tx.Get(&thread, query, id)
	return thread, err
}

func GetThreadBySlugWithoutTx(threadSlug string) (tm.Thread, error) {
	conn := utills.GetConnection()
	var thread tm.Thread
	query := `SELECT * FROM threads where slug=$1`
	err := conn.Get(&thread, query, threadSlug)
	return thread, err
}

func GetThreadByIDWithoutTx(id int) (tm.Thread, error) {
	conn := utills.GetConnection()
	var thread tm.Thread
	query := `SELECT * FROM threads where id=$1`
	err := conn.Get(&thread, query, id)
	return thread, err
}
