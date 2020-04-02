package delivery

import (
	"fmt"
	"net/http"
)

func ForumCreate(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	fmt.Println("test")
	w.WriteHeader(http.StatusOK)
}
