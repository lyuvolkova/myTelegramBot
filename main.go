package main

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/robfig/cron/v3"
	"log"
	"os"
	"strconv"
	"strings"
)

var secondParser = cron.NewParser(cron.Second | cron.Minute | cron.Hour | cron.Dom | cron.Month | cron.DowOptional | cron.Descriptor)

func getUserId(str string) (map[int64]struct{}, error) {
	parts := strings.Split(str, ",")
	users := make(map[int64]struct{})
	for _, val := range parts {
		userId, err := strconv.ParseInt(val, 10, 64)
		if err != nil {
			return nil, err
		}
		users[userId] = struct{}{}
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
	myCron := cron.New(cron.WithParser(secondParser), cron.WithChain())
	myCron.Start()
	defer myCron.Stop()
	myCron.AddFunc("0 30 18 18 * ?", func() {
		for i, _ := range users {
			_, err = bot.Send(tgbotapi.NewMessage(i, "I need data!"))
			if err != nil {
				log.Println(fmt.Errorf("cron: %w", err))
			}
		}

	})
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message != nil { // If we got a message
			if _, exists := users[update.Message.Chat.ID]; !exists { //check the sender of a message
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
