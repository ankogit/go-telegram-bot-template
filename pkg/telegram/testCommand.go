package telegram

import tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

func (b *Bot) SendTestMessage(chatId int64) error {
	msg := tgbotapi.NewMessage(chatId, "Test message")
	msg.ParseMode = "MARKDOWN"
	msg.DisableNotification = true
	b.bot.Send(msg)
	return nil
}
