package bot

import (
	amqp "github.com/rabbitmq/amqp091-go"
)

type RabbitMQ struct {
	Conn  *amqp.Connection
	Chann *amqp.Channel
	Chans Chans
}

type Chans struct {
	UserActivationChan <-chan amqp.Delivery
}
