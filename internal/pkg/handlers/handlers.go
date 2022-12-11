package handlers

import (
	"fmt"
	"github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
	"strings"
)

type handler struct {
	water waterService
	users map[int64]struct{}
	bot   *tgbotapi.BotAPI
}

func New(water waterService, users map[int64]struct{}, bot *tgbotapi.BotAPI) *handler {
	return &handler{water: water, users: users, bot: bot}
}

func (h *handler) HandleMsg(update tgbotapi.Update) error {
	//if update.Message != nil {
	//	return nil
	//}
	//fmt.Println("HERE")
	if _, exists := h.users[update.Message.Chat.ID]; !exists { //check the sender of a message
		return nil
	}

	text := update.Message.Text

	log.Printf("[%s] %s\n", update.Message.From.UserName, text)
	textMsg := "I don't understand you !"
	switch {
	case text == "Hello":
		textMsg = "Hello, " + update.Message.Chat.FirstName + "!"
		break
	case text == "Bye":
		textMsg = "Bye bye!"
		break
	case text == "/cold_data":
		fmt.Println("OK1")
		return h.getColdData(update)
	case text == "/hot_data":
		return h.getHotData(update)
	case strings.HasPrefix(text, "/cold_data "):
		return h.saveColdData(update)
	case strings.HasPrefix(text, "/hot_data "):
		return h.saveHotData(update)
	}
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, textMsg)
	//msg.ReplyToMessageID = update.Message.MessageID
	_, err := h.bot.Send(msg)
	return err
}
