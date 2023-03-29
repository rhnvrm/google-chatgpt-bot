package server

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"github.com/rhnvrm/google-chatgpt-bot/pkg/gchat"
)

func (app *App) HandleRoot(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	log.Println("Received request: ", r.Body)

	bearerToken := r.Header.Get("Authorization")
	token := strings.Split(bearerToken, " ")
	if len(token) != 2 || !gchat.VerifyJWT(app.cfg.BotAppID, token[1]) {
		http.Error(w, "Unauthorized", http.StatusForbidden)
		return
	}

	message := &gchat.Event{}
	err := json.NewDecoder(r.Body).Decode(&message)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		log.Println(err)
		return
	}

	// If the bot is added to a space, respond with a welcome message.
	if message.Type == "ADDED_TO_SPACE" {
		response := gchat.Response{Text: "Thanks for adding me!"}
		json.NewEncoder(w).Encode(response)
		return
	}

	// If the message type is not a message, ignore it.
	if message.Type != "MESSAGE" {
		response := gchat.Response{Text: "Sorry, I didn't understand your message."}
		json.NewEncoder(w).Encode(response)
		return
	}

	// Cleanup the prompt message. Argument text has a leading space usually.
	prompt := strings.TrimSpace(message.Message.ArgumentText)

	// Send the prompt to OpenAI and get a response.
	response, err := app.cfg.OpenAI.Respond(prompt, nil)
	if err != nil {
		log.Println(err)
		response := gchat.Response{Text: "Sorry, I didn't understand your message."}
		json.NewEncoder(w).Encode(response)
		return
	}

	out := gchat.Response{Text: response}
	if err := json.NewEncoder(w).Encode(out); err != nil {
		log.Println(err)
		response := gchat.Response{Text: "Sorry, I didn't understand your message."}
		json.NewEncoder(w).Encode(response)
		return
	}
}
