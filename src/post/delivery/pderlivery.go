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
