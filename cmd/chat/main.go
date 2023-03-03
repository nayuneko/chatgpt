package main

import (
	"bufio"
	"chatgpt/ai"
	"context"
	"fmt"
	"io"
	"log"
	"os"
)

func main() {
	apiToken := os.Getenv("OPENAI_API_TOKEN")
	if apiToken == "" {
		log.Fatal("環境変数OPENAI_API_TOKENが未設定")
	}
	c := ai.NewChat(apiToken)
	if err := c.SetSettingsTextFromFile("system/arisa.txt"); err != nil {
		log.Fatal(err)
	}
	ctx := context.Background()
	r := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("> ")
		line, err := r.ReadString('\n')
		if err == io.EOF {
			break
		} else if err != nil {
			log.Fatal(err)
		}
		respText, err := c.Completion(ctx, line)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(respText)
	}
}
