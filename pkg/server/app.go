package server

import (
	"log"
	"net/http"

	"github.com/rhnvrm/google-chatgpt-bot/pkg/openai"
)

type Config struct {
	OpenAIKey string
	BotAppID  string

	Address string

	OpenAI *openai.Client
}

type App struct {
	cfg Config
}

func New(cfg Config) *App {
	if cfg.Address == "" {
		cfg.Address = ":1234"
	}

	return &App{
		cfg: cfg,
	}
}

func (app *App) Run() error {
	http.HandleFunc("/", app.HandleRoot)

	log.Println("Starting server on ", app.cfg.Address)

	return http.ListenAndServe(app.cfg.Address, nil)
}
