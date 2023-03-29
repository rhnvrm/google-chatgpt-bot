package openai

import (
	"context"
	"fmt"

	"github.com/sashabaranov/go-openai"
)

type Client struct {
	client *openai.Client
	cache  map[string][]openai.ChatCompletionMessage
}

func NewClient(apiKey string) *Client {
	c := openai.NewClient(apiKey)
	return &Client{
		client: c,
		cache:  map[string][]openai.ChatCompletionMessage{},
	}
}

func (c *Client) Respond(interactionKey string, prompt string, prePrompt []openai.ChatCompletionMessage) (string, error) {
	if prePrompt == nil {
		prePrompt = []openai.ChatCompletionMessage{}
	}

	// check if we have a cached response
	if msgs, ok := c.cache[interactionKey]; ok {
		prePrompt = append(prePrompt, msgs...)
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

	// cache the response
	msgsRecv := []openai.ChatCompletionMessage{}
	for _, msg := range resp.Choices {
		msgsRecv = append(msgs, msg.Message)
	}

	c.cache[interactionKey] = append(c.cache[interactionKey], msgsRecv...)

	// if cache size is greater than 10, remove the oldest message
	if len(c.cache[interactionKey]) > 10 {
		c.cache[interactionKey] = c.cache[interactionKey][1:]
	}

	return resp.Choices[0].Message.Content, nil
}
