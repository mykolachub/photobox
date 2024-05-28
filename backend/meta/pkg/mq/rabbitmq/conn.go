package mq

import (
	"fmt"

	amqp "github.com/rabbitmq/amqp091-go"
)

type RabbitMQConfig struct {
	Host string
	Port string
	User string
}

func InitRabbitMQConnection(cfg RabbitMQConfig) (*amqp.Connection, error) {
	return amqp.Dial(fmt.Sprintf("amqp://%v@%v:%v/", cfg.User, cfg.Host, cfg.Port))
}
