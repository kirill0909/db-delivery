package usecase

import (
	"db-delivery/config"
	"db-delivery/internal/bot"
	models "db-delivery/internal/models/bot"
	"db-delivery/pkg/logger"
	// amqp "github.com/rabbitmq/amqp091-go"
)

type BotUC struct {
	log      *logger.Logger
	cfg      *config.Config
	pgRepo   bot.PgRepo
	rabbitMQ models.RabbitMQ
}

func NewBotUC(
	log *logger.Logger,
	cfg *config.Config,
	pgRepo bot.PgRepo,
	rabbitMQ models.RabbitMQ,
) bot.Usecase {
	return &BotUC{
		log:      log,
		cfg:      cfg,
		pgRepo:   pgRepo,
		rabbitMQ: rabbitMQ,
	}
}

func (u *BotUC) Consume() error {
	for {
		select {
		case msg, ok := <-u.rabbitMQ.Chans.UserActivationChan:
			if ok {
				u.log.Infof("Message: %s", string(msg.Body))
			}
		}
	}
}
