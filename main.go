package main

import (
	"log"
	"os"

	"github.com/rhnvrm/google-chatgpt-bot/pkg/openai"
	"github.com/rhnvrm/google-chatgpt-bot/pkg/server"
)

var (
	apiOpenAIKey = os.Getenv("OPENAI_SECRET_KEY")
	botAppID     = os.Getenv("BOT_APP_ID")
	address      = os.Getenv("ADDRESS")
)

func main() {
	cfg := server.Config{
		OpenAIKey: apiOpenAIKey,
		BotAppID:  botAppID,
		Address:   address,
		OpenAI:    openai.NewClient(apiOpenAIKey),
	}

	app := server.New(cfg)

	log.Fatal(app.Run())
}
