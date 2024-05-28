package mq

import (
	"context"
	"net/http"

	amqp "github.com/rabbitmq/amqp091-go"
)

type RabbitMQ struct {
	ch *amqp.Channel
}

func InitRabbitMQ(ch *amqp.Channel) *RabbitMQ {
	return &RabbitMQ{ch}
}

func (m *RabbitMQ) Publish(ctx context.Context, name string, body []byte) error {
	q, err := m.ch.QueueDeclare(
		name,  // name
		true,  // durable
		false, // delete when unused
		false, // exclusive
		false, // no-wait
		nil,   // arguments
	)
	if err != nil {
		return err
	}

	contentType := http.DetectContentType(body)
	return m.ch.PublishWithContext(ctx,
		"",
		q.Name,
		false,
		false,
		amqp.Publishing{
			DeliveryMode: amqp.Persistent,
			ContentType:  contentType,
			Body:         []byte(body),
		})

}
