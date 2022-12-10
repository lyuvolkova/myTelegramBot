package handlers

import "github.com/go-telegram-bot-api/telegram-bot-api/v5"

func (h *handler) getHotData(update tgbotapi.Update) error {
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Hot water: "+h.water.GetPrevData(5))
	_, err := h.bot.Send(msg)
	return err
}
