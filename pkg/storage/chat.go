package storage

import "github.com/ankogit/go-telegram-bot-template/pkg/models"

type ChatRepository interface {
	Save(data models.Chat) error
	Get(chatId int64) (models.Chat, error)
	Delete(chat models.Chat) error
	All() ([]models.Chat, error)
}
