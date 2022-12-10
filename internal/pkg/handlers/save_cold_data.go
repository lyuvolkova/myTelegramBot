package handlers

import tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

func (h *handler) saveColdData(update tgbotapi.Update) error {
	parTable := [2]string{"A", "B"}
	i := 0
	text := update.Message.Text
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, h.water.WriteColdData(text, parTable, i))
	_, err := h.bot.Send(msg)
	return err
}
