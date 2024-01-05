package models

import (
	bot "db-delivery/internal/models/bot"
	"github.com/jmoiron/sqlx"
	"github.com/kirill0909/logger"
)

type Dependencies struct {
	PgDB     *sqlx.DB
	Logger   *logger.Logger
	RabbitMQ bot.RabbitMQ
}
