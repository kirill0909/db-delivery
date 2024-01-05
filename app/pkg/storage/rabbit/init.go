package rabbit

import (
	"db-delivery/config"
	models "db-delivery/internal/models/bot"
	"fmt"
	"github.com/pkg/errors"
	amqp "github.com/rabbitmq/amqp091-go"
)

func InitRabbit(cfg *config.Config) (models.RabbitMQ, error) {
	url := fmt.Sprintf("%s:%s/", cfg.RabbitMQ.Host, cfg.RabbitMQ.Port)
	conn, err := amqp.Dial(url)
	if err != nil {
		err = errors.Wrap(err, "rabbit.InitRabbit. failed to connect to rabbitMQ")
		return models.RabbitMQ{}, err
	}

	ch, err := conn.Channel()
	if err != nil {
		err = errors.Wrap(err, "rabbit.InitRabbit. failed to open channel")
		return models.RabbitMQ{}, err
	}

	userActivationQueue, err := ch.QueueDeclare(
		cfg.RabbitMQ.Queues.UserActivationQueue, // name
		false,                                   // durable
		false,                                   // delete when unused
		false,                                   // exclusive
		false,                                   // no-wait
		nil,                                     // arguments
	)
	if err != nil {
		err = errors.Wrap(err, "rabbit.InitRabbit. Failed to declare queue ")
		return models.RabbitMQ{}, err
	}

	userActivationChan, err := ch.Consume(
		userActivationQueue.Name, // queue
		"",                       // consumer
		true,                     // auto-ack
		false,                    // exclusive
		false,                    // no-local
		false,                    // no-wait
		nil,                      // args
	)
	if err != nil {
		err = errors.Wrap(err, "rabbit.InitRabbit. failed to register a consumer")
		return models.RabbitMQ{}, err
	}

	return models.RabbitMQ{
		Conn:  conn,
		Chann: ch,
		Chans: models.Chans{
			UserActivationChan: userActivationChan,
		},
	}, nil
}
