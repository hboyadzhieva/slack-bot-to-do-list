package main

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/hboyadzhieva/slack-bot-to-do-list/mysql"
	"github.com/hboyadzhieva/slack-bot-to-do-list/tododo"
	"github.com/nlopes/slack"
	"log"
	"net/http"
	"os"
)

const (
	port     = ":80"
	dialect  = "mysql"
	dsn      = "myuser:mypassword@tcp(127.0.0.1:3306)/slack"
	idleConn = 10
	maxConn  = 10
)

var db *sql.DB
var commandHandler tododo.CommandHandlerInterface
var slackVerToken string

func main() {

	token, exists := os.LookupEnv("SLACK_VERIFICATION_TOKEN")
	if !exists {
		log.Fatalf("Slack verification token not set in environment")
	}
	slackVerToken = token

	db, err := sql.Open(dialect, dsn)
	if err != nil {
		log.Fatalf("Can't open DB: %s", err)
	}

	db.SetMaxIdleConns(idleConn)
	db.SetMaxOpenConns(maxConn)

	err = db.Ping()
	if err != nil {
		log.Fatalf("Ping DB error: %s", err)
	}

	defer db.Close()

	commandHandler = &tododo.CommandHandler{
		Repository: &mysql.TaskRepository{DB: db},
	}

	go http.HandleFunc("/tododo", requestHandler)
	fmt.Println("[INFO] Server listening")
	log.Fatal(http.ListenAndServe(port, nil))
}

func requestHandler(w http.ResponseWriter, r *http.Request) {
	s, err := slack.SlashCommandParse(r)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if !s.ValidateToken(slackVerToken) {
		fmt.Printf("Cant validate token %s", slackVerToken)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	response, err := commandHandler.HandleCommand(&s)
	if err != nil {
		fmt.Printf("Error handling command: %s, %s", s.Command, err)
		w.WriteHeader(http.StatusInternalServerError)
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Write(response)

}
