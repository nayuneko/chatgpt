package ai

import (
	"context"
	gogpt "github.com/sashabaranov/go-gpt3"
	"io"
	"os"
)

type Chat struct {
	client   *gogpt.Client
	messages []gogpt.ChatCompletionMessage
}

func NewChat(apiToken string) *Chat {
	return &Chat{
		client: gogpt.NewClient(apiToken),
	}

}

func (c *Chat) SetSettingsText(settingsText string) {
	if len(c.messages) != 0 {
		return
	}
	c.messages = append(c.messages, gogpt.ChatCompletionMessage{Role: "system", Content: settingsText})
}

func (c *Chat) SetSettingsTextFromFile(filename string) error {
	f, err := os.Open(filename)
	if err != nil {
		return err
	}
	b, err := io.ReadAll(f)
	if err != nil {
		return err
	}
	c.SetSettingsText(string(b))
	return nil
}

func (c *Chat) Completion(ctx context.Context, newMessageText string) (string, error) {
	c.messages = append(c.messages, gogpt.ChatCompletionMessage{Role: "user", Content: newMessageText})
	req := gogpt.ChatCompletionRequest{
		Model:    gogpt.GPT3Dot5Turbo,
		Messages: c.messages,
	}
	resp, err := c.client.CreateChatCompletion(ctx, req)
	if err != nil {
		return "", err
	}
	respText := resp.Choices[0].Message.Content
	c.messages = append(c.messages, gogpt.ChatCompletionMessage{Role: "assistant", Content: respText})
	return respText, nil
}
