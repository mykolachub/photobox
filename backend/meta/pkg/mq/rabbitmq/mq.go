package mq

import (
	amqp "github.com/rabbitmq/amqp091-go"
)

type RabbitMQ struct {
	ch *amqp.Channel
}

func InitRabbitMQ(ch *amqp.Channel) *RabbitMQ {
	return &RabbitMQ{ch}
}

func (m *RabbitMQ) Consume(name string, callback func(amqp.Delivery) error) error {
	q, err := m.ch.QueueDeclare(
		"meta_upload", // name
		true,          // durable
		false,         // delete when unused
		false,         // exclusive
		false,         // no-wait
		nil,           // arguments
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
			err := callback(d)
			if err != nil {
				continue
			}
			// var req proto.UpdateMetaRequest
			// err := pb.Unmarshal(d.Body, &req)
			// if err != nil {
			// 	log.Print("failed to unmarshal body")
			// 	continue
			// }
			// log.Printf("Received a message: %v", req.FileName)

			d.Ack(false)
		}
	}()

	<-forever
	return nil
}
