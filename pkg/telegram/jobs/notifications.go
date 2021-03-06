package jobs

type TelegramNotifications struct {
	Bot TelegramBot
}

type TelegramBot interface {
	SendTestMessage(chatId int64) error
}

func NewTelegramNotifications(bot TelegramBot) *TelegramNotifications {
	return &TelegramNotifications{Bot: bot}
}

func (t *TelegramNotifications) NotifyStats(chatId int64) error {
	return t.Bot.SendTestMessage(chatId)
}
