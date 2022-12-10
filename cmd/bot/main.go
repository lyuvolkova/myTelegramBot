package main

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/lyuvolkova/myTelegramBot/internal/pkg/handlers"
	"github.com/lyuvolkova/myTelegramBot/internal/pkg/water"
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
	//coldData := ""
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
	myCron.AddFunc("0 00 19 18 * ?", func() {
		for i, _ := range users {
			_, err = bot.Send(tgbotapi.NewMessage(i, "Hello\nhow are u?\nI NEED DATA!"))
			if err != nil {
				log.Println(fmt.Errorf("cron: %w", err))
			}
		}

	})
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)

	waterService := water.NewWaterService()
	handler := handlers.New(waterService, users, bot)
	for update := range updates {
		if err := handler.HandleMsg(update); err != nil {
			log.Println(err)
		}

	}
}
