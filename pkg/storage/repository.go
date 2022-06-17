package storage

import (
	"github.com/ankogit/go-telegram-bot-template/pkg/storage/stormDB"
	"github.com/asdine/storm/v3"
	"github.com/jmoiron/sqlx"
)

type Repositories struct {
	Chats ChatRepository
}

func NewRepositories(db *storm.DB, dbPostgres *sqlx.DB) *Repositories {
	return &Repositories{
		Chats: stormDB.NewChatRepository(db),
	}
}
