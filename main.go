package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	oidc "github.com/coreos/go-oidc"
	openai "github.com/sashabaranov/go-openai"
	chat "google.golang.org/api/chat/v1"
)

const (
	jwtURL     string = "https://www.googleapis.com/service_accounts/v1/jwk/"
	chatIssuer string = "chat@system.gserviceaccount.com"
)

var (
	apiOpenAIKey = os.Getenv("OPENAI_SECRET_KEY")
	botAppID     = os.Getenv("BOT_APP_ID")
)

func genResponse(content string) (string, error) {
	client := openai.NewClient(apiOpenAIKey)
	resp, err := client.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model: openai.GPT3Dot5Turbo,
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleUser,
					Content: content,
				},
			},
		},
	)

	if err != nil {
		return "", fmt.Errorf("ChatCompletion error: %v\n", err)
	}

	return resp.Choices[0].Message.Content, nil
}

type Message struct {
	Text string `json:"text"`
}

func main() {
	http.HandleFunc("/", handle)
	fmt.Println("Starting server on port 1234...")
	log.Fatal(http.ListenAndServe(":1234", nil))
}

func handle(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	log.Println("Received request: ", r.Body)

	// "Authorization"
	bearerToken := r.Header.Get("Authorization")
	// fmt.Println(bearerToken)
	token := strings.Split(bearerToken, " ")
	if len(token) != 2 || !verifyJWT(botAppID, token[1]) {
		http.Error(w, "Unauthorized", http.StatusForbidden)
		return
	}

	message := &chat.DeprecatedEvent{}
	err := json.NewDecoder(r.Body).Decode(&message)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		log.Println(err)
		return
	}

	content := message.Message.Text

	log.Println("Received message: ", content)

	response, err := genResponse(content)
	if err != nil {
		log.Println(err)
		response := Message{Text: "Sorry, I didn't understand your message."}
		json.NewEncoder(w).Encode(response)
		return
	}

	out := Message{Text: response}
	if err := json.NewEncoder(w).Encode(out); err != nil {
		log.Println(err)
		response := Message{Text: "Sorry, I didn't understand your message."}
		json.NewEncoder(w).Encode(response)
		return
	}
}

func verifyJWT(audience, token string) bool {
	ctx := context.Background()
	keySet := oidc.NewRemoteKeySet(ctx, jwtURL+chatIssuer)

	configLocal := &oidc.Config{
		SkipClientIDCheck: false,
		ClientID:          audience,
	}
	newVerifier := oidc.NewVerifier(chatIssuer, keySet, configLocal)
	test, err := newVerifier.Verify(ctx, token)
	if err != nil {
		log.Println("Audience doesnt match")
		return false
	}
	if len(test.Audience) == 1 {
		log.Println("Audience matches")
	}

	return true
}
