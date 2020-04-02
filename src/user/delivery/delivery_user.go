package delivery

import (
	"../../utills"
	umodel "../models"
	urepo "../repository"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gorilla/mux"
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
	nickname := mux.Vars(r)["nickname"]
	user, err := fetchUser(r)
	if err != nil {
		utills.SendServerError("can Json Data", http.StatusConflict, w)
		return
	}
	user.Nickname = nickname
	err = urepo.SaveUser(user)
	fmt.Println(err)
	if err != nil {
		utills.SendServerError("user already ", http.StatusConflict, w)
		return
	}

	utills.SendOKAnswer(user, w)

}

func UserGetOne(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	nickname := mux.Vars(r)["nickname"]
	user, err := urepo.GetUserByNickname(nickname)
	fmt.Println(err)
	if err != nil {
		utills.SendServerError("cant find user with nickname "+nickname, http.StatusUnauthorized, w)
		return
	}
	utills.SendOKAnswer(user, w)
}
