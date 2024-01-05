package usecase

import (
	"context"
	"db-delivery/config"
	"db-delivery/internal/bot"
	models "db-delivery/internal/models/bot"
	"db-delivery/pkg/logger"
	"encoding/json"
	"github.com/pkg/errors"
	amqp "github.com/rabbitmq/amqp091-go"
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
				if err := u.handleUserActivationCase(msg); err != nil {
					u.log.Errorf(err.Error())
				}
			}
		}
	}
}

func (u *BotUC) handleUserActivationCase(msg amqp.Delivery) error {
	u.log.Infof("read message: (%+v) from queue(%s)", string(msg.Body), u.cfg.RabbitMQ.Queues.UserActivationQueue)
	var userActivationParams models.UserActivationParams
	if err := json.Unmarshal(msg.Body, &userActivationParams); err != nil {
		return errors.Wrapf(err, "BotUC.handleUserActivationCase.Unmrashal. param(%+v)", string(msg.Body))
	}

	ctx := context.Background()
	if err := u.pgRepo.UserActivation(ctx, userActivationParams); err != nil {
		return err
	}

	return nil
}
