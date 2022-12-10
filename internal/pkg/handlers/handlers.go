package handlers

import (
	"github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
	"strings"
)

func HandlersMsg(water waterService, update tgbotapi.Update, users map[int64]struct{}, bot *tgbotapi.BotAPI) error {
	if update.Message != nil {
		return nil
	}

	if _, exists := users[update.Message.Chat.ID]; !exists { //check the sender of a message
		return nil
	}

	text := update.Message.Text
	log.Printf("[%s] %s\n", update.Message.From.UserName, text)
	textMsg := "I don't understand you !"
	if text == "Hello" {
		textMsg = "Hello, " + update.Message.Chat.FirstName + "!"
	} else if text == "Bye" {
		textMsg = "Bye bye!"
	} else if text == "/cold_data" {
		textMsg = "Cold water: " + water.GetPrevData(1)
	} else if text == "/hot_data" {
		textMsg = "Hot water: " + water.GetPrevData(5)
	} else if strings.HasPrefix(text, "/cold_data ") {
		parTable := [2]string{"A", "B"}
		i := 0
		textMsg = water.WriteColdData(text, parTable, i)
	} else if strings.HasPrefix(text, "/hot_data ") {
		parTable := [2]string{"E", "F"}
		i := 4
		textMsg = water.WriteColdData(text, parTable, i)
	}
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, textMsg)
	//msg.ReplyToMessageID = update.Message.MessageID
	_, err := bot.Send(msg)
	return err
}
