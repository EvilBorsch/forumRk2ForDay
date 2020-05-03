package post

import (
	"errors"
	"fmt"
	"github.com/jmoiron/sqlx"
	pmodel "go-server-server-generated/src/post/models"
	tmodel "go-server-server-generated/src/thread/models"
	trepo "go-server-server-generated/src/thread/repository"
	"go-server-server-generated/src/user/repository"
	"go-server-server-generated/src/utills"
	"strconv"
	"time"
)

func IsDigit(v string) bool {
	if _, err := strconv.Atoi(v); err == nil {
		return true
	}
	return false
}

func UpdateForumPostsCountByThread(tx *sqlx.Tx, thread tmodel.Thread, incValue int) error {
	forumSlug := thread.Forum
	query := `UPDATE forum SET posts=posts+$1 where slug=$2`
	_, err := tx.Exec(query, incValue, forumSlug)
	return err
}

func CheckIfParentPostsInSameThread(tx *sqlx.Tx, post pmodel.Post) bool {
	if post.Parent == nil {
		return true
	}
	parentPost, err := GetPostById(tx, post.Parent)
	if err != nil {
		return false
	}
	return parentPost.Thread == post.Thread
}

func GetPostById(tx *sqlx.Tx, id *int) (pmodel.Post, error) {
	query := `Select * from posts where id=$1`
	var post pmodel.Post
	err := tx.Get(&post, query, id)
	return post, err
}

func AddNewPosts(posts []pmodel.Post, thr tmodel.Thread) ([]pmodel.Post, error) {
	timeCreated := time.Now().UTC()
	query := `INSERT INTO posts (author, created, forum, isedited, message, parent, thread) VALUES ($1,$2,$3,$4,$5,NULLIF($6,0),$7) returning *`
	conn := utills.GetConnection()
	tx := conn.MustBegin()
	defer tx.Commit()
	var postList []pmodel.Post
	var err error
	thread := thr

	for _, post := range posts {
		post.Created = timeCreated
		post.Forum = thread.Forum
		post.Thread = thread.Id
		post.IsEdited = false

		ok := checkIfAuthorExist(tx, post)
		if !ok {
			errMsg := "Can't find post author by nickname: " + post.Author
			return nil, errors.New(errMsg)
		}
		ok = CheckIfParentPostsInSameThread(tx, post)
		if !ok {
			return nil, errors.New("no parent")
		}

		var newPost pmodel.Post
		err := tx.Get(&newPost, query, post.Author, post.Created, post.Forum, post.IsEdited, post.Message, post.Parent, post.Thread)
		fmt.Println("E:", err, newPost)
		fmt.Println("new post: ", newPost)
		newPost.Thread = post.Thread //COSTIL todo
		postList = append(postList, newPost)

	}

	err = UpdateForumPostsCountByThread(tx, thread, len(postList))
	fmt.Println(postList)
	return postList, err
}

func checkIfAuthorExist(tx *sqlx.Tx, post pmodel.Post) bool {
	_, err := repository.GetUserByNickname(post.Author)
	if err != nil {
		return false
	}
	return true
}

func GetPostsWithFlatSort(slug_or_id string, limit int, sinceID int) ([]pmodel.Post, error) {
	tx := utills.StartTransaction()
	defer utills.EndTransaction(tx)
	id, isDig := utills.IsDigit(slug_or_id)
	if isDig {
		return GetPostsWithFlatSortById(tx, id, limit, sinceID)
	} else {
		return GetPostsWithFlatSortBySlug(tx, slug_or_id, limit, sinceID)
	}
}

func GetPostsWithFlatSortBySlug(tx *sqlx.Tx, slug string, limit int, sinceId int) ([]pmodel.Post, error) {
	thread, err := trepo.GetThreadBySlug(tx, slug)
	if err != nil {
		return nil, err
	}
	return GetPostsWithFlatSortById(tx, thread.Id, limit, sinceId)

}

func GetPostsWithFlatSortById(tx *sqlx.Tx, id int, limit int, sinceID int) ([]pmodel.Post, error) {
	query := `Select * from posts where thread=$1 and id>$2 order by created,id LIMIT $3`
	var posts []pmodel.Post
	err := tx.Select(&posts, query, id, sinceID, limit)
	fmt.Println(err)
	return posts, err
}

func UpdatePost(tx *sqlx.Tx, postId int, UpdatedPost pmodel.Post) (pmodel.Post, error) {
	prevPost, _ := GetPostById(tx, &postId)
	if prevPost.Message == UpdatedPost.Message || UpdatedPost.Message == "" {
		return prevPost, nil
	}
	query := `UPDATE posts SET message=$1,isEdited=true where id=$2 returning *`
	var post pmodel.Post
	err := tx.Get(&post, query, UpdatedPost.Message, postId)
	return post, err

}
