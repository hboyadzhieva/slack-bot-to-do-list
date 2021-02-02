package main

import (
	"fmt"
	"github.com/joho/godotenv"
	"github.com/nlopes/slack"
	"log"
	"net/http"
	"os"
)

func main() {
	err := godotenv.Load("keys.env")
	if err != nil {
		log.Fatal("Error loading environment")
	}

	http.HandleFunc("/", slashCommandHandler)
	fmt.Println("[INFO] Server listening")
	log.Fatal(http.ListenAndServe(":80", nil))
}

func slashCommandHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("I heard ya")
	s, err := slack.SlashCommandParse(r)
	fmt.Printf("%+v\n", s)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if !s.ValidateToken(os.Getenv("SLACK_VERIFICATION_TOKEN")) {
		fmt.Println("Slack couldnt ValidateToken")
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	switch s.Command {
	case "/tdcreatelist":
		fmt.Println("Recognized your command")
		params := &slack.Msg{Text: s.Text}
		response := fmt.Sprintf("You want to create a list %v", params.Text)
		w.Write([]byte(response))
	default:
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
