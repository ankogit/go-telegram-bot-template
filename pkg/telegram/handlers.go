package telegram

import (
	"github.com/ankogit/go-telegram-bot-template/pkg/models"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
	"regexp"
	"strings"
)

func (b *Bot) handleInlineQuery(query *tgbotapi.InlineQuery) {
	var resources []interface{}

	resources = append(resources,
		tgbotapi.InlineQueryResultArticle{
			Type:  "article",
			ID:    query.ID,
			Title: b.messages.InlineContentTitle,
			InputMessageContent: tgbotapi.InputTextMessageContent{
				Text: "Some description"},
			Description: b.messages.InlineContentDescription})

	if _, err := b.bot.Request(tgbotapi.InlineConfig{
		InlineQueryID: query.ID,
		CacheTime:     0,
		IsPersonal:    true,
		Results:       resources}); err != nil {
		log.Println(err)
	}
}

func (b *Bot) handleCommand(message *tgbotapi.Message) error {
	switch message.Command() {
	case "start":
		b.SendWelcomeMessage(message.Chat.ID)
		return nil
	case "test":
		if err := b.handleCommandTest(message.Chat.ID); err != nil {
			return err
		}
		return nil
	default:
		return errUnknownCommand
	}
}

func (b *Bot) SendWelcomeMessage(chatId int64) {
	msg := tgbotapi.NewMessage(chatId, "Hello world! Type @"+b.bot.Self.UserName+" in message field. \nv. "+b.version+"")
	b.bot.Send(msg)
}

func (b *Bot) handleCommandTest(chatId int64) error {
	return b.SendTestMessage(chatId)
}

func (b *Bot) handleNotificationEnable(message *tgbotapi.Message) error {
	if message.Chat == nil || message.Chat.Title == "" {
		msg := tgbotapi.NewMessage(message.Chat.ID, "Notifications are allowed to conduct group chats")
		b.bot.Send(msg)
		return nil
	}
	cParts := strings.SplitAfterN(message.Text, " ", 2)
	if len(cParts) == 1 {
		msg := tgbotapi.NewMessage(message.Chat.ID, "To set notification you need to pass the cron param")
		b.bot.Send(msg)
		return nil
	}
	cronParam := cParts[1]

	var validID = regexp.MustCompile(`(@(annually|yearly|monthly|weekly|daily|hourly|reboot))|(@every (\d+(ns|us|µs|ms|s|m|h))+)|((((\d+,)+\d+|(\d+(\/|-)\d+)|\d+|\*) ?){5,7})`)
	if !validID.MatchString(cronParam) {
		msg := tgbotapi.NewMessage(message.Chat.ID, "Cron string is not correct :(")
		b.bot.Send(msg)
		return nil
	}

	var chat = models.Chat{
		ID:               message.Chat.ID,
		Title:            message.Chat.Title,
		NotificationCron: cronParam,
		EntryId:          0,
	}
	if err := b.chatRepository.Save(chat); err != nil {
		return err
	}

	entryId, err := b.cronService.SetJob(&chat, cronParam)
	if err != nil {
		return err
	}
	chat.EntryId = entryId

	if err := b.chatRepository.Save(chat); err != nil {
		return err
	}
	msg := tgbotapi.NewMessage(message.Chat.ID, "Notifications are installed for this chat 🔔")
	b.bot.Send(msg)
	return nil
}

func (b *Bot) handleNotificationDisable(message *tgbotapi.Message) error {
	if message.Chat == nil || message.Chat.Title == "" {
		msg := tgbotapi.NewMessage(message.Chat.ID, "Notifications are allowed to conduct group chats")
		b.bot.Send(msg)
		return nil
	}

	chat, _ := b.chatRepository.Get(message.Chat.ID)

	if (chat == (models.Chat{})) || (chat.NotificationCron == "" && chat.EntryId == 0) {
		msg := tgbotapi.NewMessage(message.Chat.ID, "Notifications are not already setting for this chat")
		b.bot.Send(msg)
		return nil
	}

	b.cronService.RemoveJob(&chat)
	chat.NotificationCron = ""
	chat.EntryId = 0

	if err := b.chatRepository.Save(chat); err != nil {
		return err
	}
	msg := tgbotapi.NewMessage(message.Chat.ID, "Notifications are disabled for this chat 🔕")
	b.bot.Send(msg)
	return nil
}
