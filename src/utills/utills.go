package utills

import (
	"encoding/json"
	"fmt"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/rs/zerolog/log"
	"net/http"
)

var conn *sqlx.DB

func CreateConnection() {
	connStr := fmt.Sprintf("user=%s password=%s dbname=postgres sslmode=disable port=%s",
		"postgres",
		"",
		"5432")
	conn, _ = sqlx.Open("postgres", connStr)

}

func GetConnection() *sqlx.DB {
	return conn
}

type HttpError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type HttpResponse struct {
	Data   interface{} `json:"data,omitempty"`
	Errors []HttpError `json:"errors"`
}

func SendServerError(errorMessage string, code int, w http.ResponseWriter) {
	log.Error().Msgf(errorMessage)
	w.WriteHeader(code)
	mes, _ := json.Marshal(errorMessage)
	w.Write(mes)
}

func SendOKAnswer(data interface{}, w http.ResponseWriter) {
	w.WriteHeader(http.StatusCreated)
	serializedData, err := json.Marshal(data)
	if err != nil {
		log.Error().Msgf(err.Error())
		return
	}

	_, err = w.Write(serializedData)
	if err != nil {
		message := fmt.Sprintf("HttpResponse while writing is socket: %s", err.Error())
		log.Error().Msgf(message)
		return
	}
	log.Info().Msgf("OK message sent")
}
