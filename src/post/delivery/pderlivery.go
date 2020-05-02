package post

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	pmodel "go-server-server-generated/src/post/models"
	prepo "go-server-server-generated/src/post/repository"
	"go-server-server-generated/src/utills"
	"io/ioutil"
	"net/http"
	"strconv"
)

func PostsCreate(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	posts, err := fetchPost(r)
	fmt.Println(posts, err)
	threadSlug := mux.Vars(r)["slug_or_id"]

	if len(posts) == 0 {
		utills.SendAnswerWithCode([]pmodel.Post{}, http.StatusCreated, w)
		return
	}

	newPosts, err := prepo.AddNewPosts(posts, threadSlug)
	fmt.Println(err)
	if err != nil {
		if err.Error() == "no thread" {
			utills.SendServerError("no thread", http.StatusNotFound, w)
			return
		}
		if err.Error()[0] == 'C' {
			utills.SendServerError(err.Error(), 404, w)
			return
		} else {
			utills.SendServerError("no parent", http.StatusConflict, w)
			return
		}
	}

	utills.SendAnswerWithCode(newPosts, http.StatusCreated, w)

}

func fetchPost(r *http.Request) ([]pmodel.Post, error) {
	defer r.Body.Close()
	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return []pmodel.Post{}, err
	}
	var post []pmodel.Post
	err = json.Unmarshal(data, &post)
	return post, err
}

func GetPosts(w http.ResponseWriter, r *http.Request) {
	var limit int
	limitStr := r.FormValue("limit")
	if limitStr == "" {
		limit = 100
	} else {
		limit, _ = strconv.Atoi(limitStr)
	}
	slug_or_id := mux.Vars(r)["slug_or_id"]
	sort := r.FormValue("sort")
	fmt.Println(limit, sort)
	if sort == "" || sort == "flat" {
		posts, err := prepo.GetPostsWithFlatSort(slug_or_id, limit)
		if err != nil {
			utills.SendServerError("posts not found", 404, w)
			return
		}
		utills.SendOKAnswer(posts, w)
		return
	}

}
