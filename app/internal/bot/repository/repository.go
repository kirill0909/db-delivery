package repository

import (
	"db-delivery/internal/bot"
	"github.com/jmoiron/sqlx"
)

type BotPGRepo struct {
	db *sqlx.DB
}

func NewBotPGRepo(db *sqlx.DB) bot.PgRepo {
	return &BotPGRepo{db: db}
}
