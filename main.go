package main

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/hboyadzhieva/slack-bot-to-do-list/mysql"
	"github.com/hboyadzhieva/slack-bot-to-do-list/tododo"
	"github.com/joho/godotenv"
	"github.com/nlopes/slack"
	"log"
	"net/http"
	"os"
)

const (
	port     = ":80"
	dialect  = "mysql"
	dsn      = "root:pass@tcp(127.0.0.1:3306)/slack"
	idleConn = 10
	maxConn  = 10
)

var db *sql.DB

func main() {
	err := godotenv.Load(os.Getenv("GOPATH") + string(os.PathSeparator) + "keys.env")
	if err != nil {
		log.Fatal("Error loading environment", err)
	}

	db, err = sql.Open(dialect, dsn)
	if err != nil {
		log.Fatalf("Can't open DB: %s", err)
	}

	err = db.Ping()
	if err != nil {
		log.Fatalf("Ping DB error: %s", err)
	}

	db.SetMaxIdleConns(idleConn)
	db.SetMaxOpenConns(maxConn)

	defer db.Close()

	http.HandleFunc("/tododo", requestHandler)
	fmt.Println("[INFO] Server listening")
	log.Fatal(http.ListenAndServe(port, nil))
}

func requestHandler(w http.ResponseWriter, r *http.Request) {
	s, err := slack.SlashCommandParse(r)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if !s.ValidateToken(os.Getenv("SLACK_VERIFICATION_TOKEN")) {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	commandHandler := &tododo.CommandHandler{
		Repository: &mysql.TaskRepository{DB: db},
	}

	response, err := commandHandler.HandleCommand(&s)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(response))

}
