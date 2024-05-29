package mq

import (
	"context"
	"log"
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

func (m *RabbitMQ) Consume(name string, callback func(amqp.Delivery)) error {
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

	err = m.ch.Qos(
		1,     // prefetch count
		0,     // prefetch size
		false, // global
	)
	if err != nil {
		return err
	}

	msgs, err := m.ch.Consume(
		q.Name, // queue
		"",     // consumer
		false,  // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	if err != nil {
		return err
	}

	var forever chan struct{}

	go func() {
		for d := range msgs {
			callback(d)
			d.Ack(false)
		}
	}()

	log.Printf("RabbitMQ Waiting for [%v] messages...", name)
	<-forever
	return nil
}
