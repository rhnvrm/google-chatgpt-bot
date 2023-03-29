package openai

import (
	"context"
	"fmt"

	"github.com/sashabaranov/go-openai"
)

type Client struct {
	client *openai.Client
}

func NewClient(apiKey string) *Client {
	c := openai.NewClient(apiKey)
	return &Client{
		client: c,
	}
}

func (c *Client) Respond(prompt string, prePrompt []openai.ChatCompletionMessage) (string, error) {
	if prePrompt == nil {
		prePrompt = []openai.ChatCompletionMessage{}
	}

	msgs := append(prePrompt, openai.ChatCompletionMessage{
		Role:    openai.ChatMessageRoleUser,
		Content: prompt,
	})

	resp, err := c.client.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model:    openai.GPT3Dot5Turbo,
			Messages: msgs,
		},
	)

	if err != nil {
		return "", fmt.Errorf("ChatCompletion error: %v\n", err)
	}

	return resp.Choices[0].Message.Content, nil
}
