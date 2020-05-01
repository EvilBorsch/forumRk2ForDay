package post

import (
	"errors"
	"fmt"
	pmodel "go-server-server-generated/src/post/models"
	tmodel "go-server-server-generated/src/thread/models"
	trepo "go-server-server-generated/src/thread/repository"
	"go-server-server-generated/src/utills"
	"strconv"
	"time"
)

func isDigit(v string) bool {
	if _, err := strconv.Atoi(v); err == nil {
		return true
	}
	return false
}

func AddNewPosts(posts []pmodel.Post, threadSlug string) ([]pmodel.Post, error) {
	query := `INSERT INTO posts (author, created, forum, isedited, message, parent, thread) VALUES ($1,$2,$3,$4,$5,NULLIF($6,0),$7) returning *`
	conn := utills.GetConnection()
	tx := conn.MustBegin()
	var postList []pmodel.Post
	var err error
	var thread tmodel.Thread
	if isDigit(threadSlug) {
		fmt.Println("is digit")
		//todo
		id, _ := strconv.Atoi(threadSlug)
		thread, err = trepo.GetThreadByID(tx, id)
	} else {
		thread, err = trepo.GetThreadBySlug(tx, threadSlug)
	}
	if err != nil {
		return nil, errors.New("no thread")
	}
	fmt.Println(err)
	for _, post := range posts {
		post.Created = time.Now().UTC()
		post.Forum = thread.Forum
		post.Thread = thread.Id
		post.IsEdited = false
		fmt.Println(post)
		var newPost pmodel.Post
		err := tx.Get(&newPost, query, post.Author, post.Created, post.Forum, post.IsEdited, post.Message, post.Parent, post.Thread)
		fmt.Println("E:", err, newPost)
		fmt.Println("new post: ", newPost)
		newPost.Thread = post.Thread //COSTIL todo
		postList = append(postList, newPost)

	}
	tx.Commit()
	fmt.Println(postList)
	return postList, err
}
