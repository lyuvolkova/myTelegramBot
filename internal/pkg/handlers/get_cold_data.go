package handlers

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (h *handler) getColdData(update tgbotapi.Update) error {
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Cold water: "+h.water.GetPrevData(1))
	_, err := h.bot.Send(msg)
	return err
}
