package main

import (
	"bufio"
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
	coldData := "0.0"
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
			text := update.Message.Text
			log.Printf("[%s] %s\n", update.Message.From.UserName, text)
			textMsg := "I don't understand you !"
			if text == "Hello" {
				textMsg = "Hello, " + update.Message.Chat.FirstName + "!"
			} else if text == "Bye" {
				textMsg = "Bye bye!"
			} else if text == "/cold_data" {
				file, err := os.Open("coldData.txt")
				if err != nil {
					log.Fatal(err)
				}
				defer file.Close()

				scanner := bufio.NewScanner(file)
				for scanner.Scan() {
					coldData = scanner.Text()
				}
				//s := fmt.Sprintf("%v", coldData)
				if coldData == "" {
					coldData = "No previous data"
				}
				textMsg = "Cold water: " + coldData
			} else if strings.HasPrefix(text, "/cold_data ") {
				coldData = string([]rune(text)[11:])
				_, err = strconv.ParseFloat(coldData, 32)
				if err == nil {
					coldData += "\n"
					file, err := os.OpenFile("coldData.txt", os.O_APPEND|os.O_WRONLY, 0600)
					if err != nil {
						panic(err)
					}
					defer file.Close()

					if _, err = file.WriteString(coldData); err != nil {
						panic(err)
					}
					textMsg = "Ok, date saved"
				} else {
					textMsg = "Incorrect number"
				}
			}
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, textMsg)
			//msg.ReplyToMessageID = update.Message.MessageID
			bot.Send(msg)
		}
	}
}
