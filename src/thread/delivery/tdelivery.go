package delivery

import (
	"../../utills"
	tmodel "../models"
	trepo "../repository"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
)

func fetchThread(r *http.Request) (tmodel.Thread, error) {
	var badStaff = errors.New("bad json data")
	defer r.Body.Close()
	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return tmodel.Thread{}, badStaff
	}
	var thread tmodel.Thread
	err = json.Unmarshal(data, &thread)
	thread.Created = thread.Created.UTC()
	return thread, err
}

func ThreadCreate(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	thread, err := fetchThread(r)
	fmt.Println("thread: ", thread, err)
	resultThread, err := trepo.AddNew(thread)

	fmt.Println(resultThread)
	if err != nil {
		fmt.Println("err: ", err)
		if err.Error() == `pq: insert or update on table "threads" violates foreign key constraint "threads_author_fkey"` {
			utills.SendServerError("not found", http.StatusNotFound, w)
			return
		}
	}
	resultThread.Created = resultThread.Created.UTC()
	utills.SendAnswerWithCode(resultThread, http.StatusCreated, w)
}
