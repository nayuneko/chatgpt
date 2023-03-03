package main

import (
	"chatgpt/ai"
	"context"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/line/line-bot-sdk-go/linebot"
	"log"
	"os"
)

var envKeys = []string{"OPENAI_API_TOKEN", "LINE_CHANNEL_SECRET", "LINE_CHANNEL_TOKEN"}

func loadConfig() map[string]string {
	godotenv.Load(".env")
	config := map[string]string{}
	for _, k := range envKeys {
		v := os.Getenv(k)
		if v == "" {
			log.Fatalf("環境変数%sが未設定", k)
		}
		config[k] = v
	}
	return config
}

func main() {

	config := loadConfig()

	bot, err := linebot.New(
		config["LINE_CHANNEL_SECRET"],
		config["LINE_CHANNEL_TOKEN"],
	)
	if err != nil {
		log.Fatal(err)
	}

	chat := ai.NewChat(config["OPENAI_API_TOKEN"])
	if err := chat.SetSettingsTextFromFile("system/arisa.txt"); err != nil {
		log.Fatal(err)
	}

	ctx := context.Background()

	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.POST("/arisa/callback", func(c echo.Context) error {
		events, err := bot.ParseRequest(c.Request())
		if err != nil {
			if err == linebot.ErrInvalidSignature {
				log.Print(err)
				return nil
			}
			log.Fatal(err)
		}
		for _, ev := range events {
			if ev.Type != linebot.EventTypeMessage {
				continue
			}
			// TODO: ユーザ毎にセッション管理した方がよい ev.Source.UserID
			if msg, ok := ev.Message.(*linebot.TextMessage); ok {
				respText, err := chat.Completion(ctx, msg.Text)
				if err != nil {
					log.Fatal(err)
				}
				bot.ReplyMessage(ev.ReplyToken, linebot.NewTextMessage(respText)).Do()
			}
		}
		return nil
	})
	e.Start(":20069")
}
