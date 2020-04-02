package delivery

import (
	"../../utills"
	fmodel "../models"
	frepo "../repository"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gorilla/mux"
	"io/ioutil"
	"net/http"
)

var badStaff = errors.New("bad json data")

func fetchForum(r *http.Request) (fmodel.Forum, error) {
	defer r.Body.Close()
	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return fmodel.Forum{}, badStaff
	}
	var forum fmodel.Forum
	err = json.Unmarshal(data, &forum)
	return forum, err
}

func ForumCreate(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	forum, err := fetchForum(r)

	newForum, err := frepo.CreateForum(forum)
	fmt.Println(newForum, err)
	if err != nil && (err.Error() == `pq: insert or update on table "forum" violates foreign key constraint "forum_user_nickname_fkey"` || err == sql.ErrNoRows) {
		utills.SendServerError("cant find user", http.StatusNotFound, w)
		return
	}
	if err != nil && err.Error() == `pq: duplicate key value violates unique constraint "forum_slug_key"` {
		forumWithThisSlug, err := frepo.GetForumBySlug(forum.Slug)
		fmt.Println("ASDJHASKJHD: ", err)
		utills.SendAnswerWithCode(forumWithThisSlug, http.StatusConflict, w)
		return
	}
	utills.SendAnswerWithCode(newForum, http.StatusCreated, w)
}

func ForumGetOne(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	slug := mux.Vars(r)["slug"]
	forum, err := frepo.GetForumBySlug(slug)
	if err != nil {
		utills.SendServerError("forum not in", http.StatusNotFound, w)
		return
	}

	utills.SendOKAnswer(forum, w)
}
