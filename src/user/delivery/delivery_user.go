package delivery

import (
	"../../utills"
	umodel "../models"
	urepo "../repository"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
)

var badStaff = errors.New("bad data")

func fetchUser(r *http.Request) (umodel.User, error) {
	defer r.Body.Close()
	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return umodel.User{}, badStaff
	}
	var user umodel.User
	err = json.Unmarshal(data, &user)
	return user, nil

}

func UserCreate(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	user, err := fetchUser(r)
	fmt.Println(user, err)
	utills.SendOKAnswer(user, w)
	err = urepo.SaveUser(user)
	fmt.Println(err)

}
