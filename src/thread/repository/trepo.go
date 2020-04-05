package repository

import (
	frepo "../../forum/repository"
	"../../utills"
	tm "../models"
	"errors"
	"fmt"
	"github.com/jmoiron/sqlx"
	"strconv"
)

func AddNew(newThread tm.Thread) (tm.Thread, error) {
	conn := utills.GetConnection()
	query := `INSERT INTO threads (author, created, forum, message, title,slug) VALUES ($1,$2,$3,$4,$5,nullif($6,'')) returning *`
	var createdThred tm.Thread
	forum, err := frepo.GetForumBySlug(newThread.Forum) //kostil for tests
	if err != nil {
		return tm.Thread{}, errors.New("no forum")
	}
	newThread.Forum = forum.Slug

	err = conn.Get(&createdThred, query, newThread.Author, newThread.Created, newThread.Forum, newThread.Message, newThread.Title, newThread.Slug)
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

func isDigit(v string) bool {
	if _, err := strconv.Atoi(v); err == nil {
		return true
	}
	return false
}

func IncrementVoteBySlug(tx *sqlx.Tx, slug string, delta int) (tm.Thread, error) {
	query := `UPDATE threads SET votes =votes + $1 WHERE slug=$2 returning *`
	var thread tm.Thread
	err := tx.Get(&thread, query, delta, slug)
	return thread, err
}

func IncrementVoteByID(tx *sqlx.Tx, id string, delta int) (tm.Thread, error) {
	query := `UPDATE threads SET votes =votes + $1 WHERE id=$2 returning *`
	var thread tm.Thread
	err := tx.Get(&thread, query, delta, id)
	return thread, err
}

func UpdateVoteInVotes(tx *sqlx.Tx, newVote tm.Vote) error {
	query := `UPDATE votes SET voice =$1 WHERE nickname=$2`
	_, err := tx.Exec(query, newVote.Voice, newVote.Nickname)
	fmt.Println("upd: ", err)
	return err
}

func getOldVoteByAuthor(tx *sqlx.Tx, author string) (tm.Vote, error) {
	query := `Select nickname,voice from votes where nickname=$1`
	var oldVote tm.Vote
	err := tx.Get(&oldVote, query, author)
	if err != nil {
		notFound := errors.New("no vote")
		return tm.Vote{}, notFound
	}
	return oldVote, nil
}

func getOldVoteByIdAndAuthor(tx *sqlx.Tx, id string, author string) (tm.Vote, error) {
	query := `Select nickname,voice from votes where nickname=$1 and threadID=$2 `
	var oldVote tm.Vote
	err := tx.Get(&oldVote, query, author, id)
	if err != nil {
		notFound := errors.New("no vote")
		return tm.Vote{}, notFound
	}
	return oldVote, nil
}

func InsertNewVoteWithThreadId(tx *sqlx.Tx, newVote tm.Vote, slug_or_id string) {
	query := `INSERT INTO votes (threadid, threadslug, nickname, voice) VALUES ($1,$2,$3,$4)`
	tx.Exec(query, slug_or_id, nil, newVote.Nickname, newVote.Voice)
}

func InsertNewVoteWithThreadSlug(tx *sqlx.Tx, newVote tm.Vote, slug_or_id string) {

	query := `INSERT INTO votes (threadid, threadslug, nickname, voice) VALUES ($1,$2,$3,$4)`
	tx.Exec(query, nil, slug_or_id, newVote.Nickname, newVote.Voice)

}

func MakeVote(slug_or_id string, newVote tm.Vote) (tm.Thread, error) {
	AlreadyExistErr := errors.New("vote already exist")
	ThreadIsNotExistErr := errors.New("thread is not exist")
	conn := utills.GetConnection()
	tx := conn.MustBegin()
	defer tx.Commit()
	if isDigit(slug_or_id) {
		oldVote, err := getOldVoteByAuthor(tx, newVote.Nickname)
		if err != nil && err.Error() == `no vote` {
			delta := newVote.Voice
			IncThread, err := IncrementVoteByID(tx, slug_or_id, delta)
			if err != nil {
				return tm.Thread{}, ThreadIsNotExistErr
			}
			InsertNewVoteWithThreadId(tx, newVote, slug_or_id)

			return IncThread, err
		}
		if oldVote.Voice == newVote.Voice {
			return tm.Thread{}, AlreadyExistErr
		}
		delta := newVote.Voice - oldVote.Voice
		IncThread, err := IncrementVoteByID(tx, slug_or_id, delta)
		if err != nil {
			return tm.Thread{}, ThreadIsNotExistErr
		}
		UpdateVoteInVotes(tx, newVote)
		return IncThread, err
	}

	oldVote, err := getOldVoteByAuthor(tx, newVote.Nickname)
	if err != nil && err.Error() == `no vote` {
		delta := newVote.Voice
		IncThread, err := IncrementVoteBySlug(tx, slug_or_id, delta)
		if err != nil {
			return tm.Thread{}, ThreadIsNotExistErr
		}
		InsertNewVoteWithThreadSlug(tx, newVote, slug_or_id)

		return IncThread, err
	}
	if oldVote.Voice == newVote.Voice {
		return tm.Thread{}, AlreadyExistErr
	}
	delta := newVote.Voice - oldVote.Voice
	IncThread, err := IncrementVoteBySlug(tx, slug_or_id, delta)
	if err != nil {
		return tm.Thread{}, ThreadIsNotExistErr
	}
	UpdateVoteInVotes(tx, newVote)
	return IncThread, err

}
