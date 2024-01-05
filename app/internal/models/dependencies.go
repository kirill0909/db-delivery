package models

import (
	bot "db-delivery/internal/models/bot"
	"db-delivery/pkg/logger"
	"github.com/jmoiron/sqlx"
)

type Dependencies struct {
	PgDB     *sqlx.DB
	Logger   *logger.Logger
	RabbitMQ bot.RabbitMQ
}
