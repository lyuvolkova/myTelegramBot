package main

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
	"os"
	"strconv"
	"strings"
)

func getUserId(str string) (map[int64]struct{}, error) {
	parts := strings.Split(str, ",")
	users := make(map[int64]struct{})
	for _, val := range parts {
		user_id, err := strconv.ParseInt(val, 10, 64)
		if err != nil {
			return nil, err
		}
		users[user_id] = struct{}{}
	}
	return users, nil
}

func main() {
	users, err := getUserId(os.Getenv("USER_IDS"))
	if err != nil {
		log.Panic(err)
	}
	bot, err := tgbotapi.NewBotAPI(os.Getenv("TOKEN"))
	if err != nil {
		log.Panic(err)
	}
	bot.Debug = true
	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message != nil { // If we got a message
			if _, exists := users[update.Message.Chat.ID]; !exists {
				continue
			}
			log.Printf("[%s] %s\n", update.Message.From.UserName, update.Message.Text)
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "I don't understand you !")
			if update.Message.Text == "Hello" {
				msg = tgbotapi.NewMessage(update.Message.Chat.ID, "Hello!")
			} else if update.Message.Text == "Bye" {
				msg = tgbotapi.NewMessage(update.Message.Chat.ID, "Bye bye!")
			}
			msg.ReplyToMessageID = update.Message.MessageID
			bot.Send(msg)
		}
	}
}
