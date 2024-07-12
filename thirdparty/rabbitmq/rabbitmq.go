package rabbitmq

import (
	"context"
	"encoding/json"
	"log"
	"twitch_chat_analysis/pkg/models"

	amqp "github.com/rabbitmq/amqp091-go"
)

type RabbitMQ struct {
	Channel *amqp.Channel
	Queue   amqp.Queue
}

func NewRabbitMQInstance() RabbitMQ {
	conn, err := amqp.Dial("amqp://user:password@localhost:7001/")
	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ: %s", err)
	}

	ch, err := conn.Channel()
	if err != nil {
		conn.Close()
		log.Fatalf("Failed to open a channel: %s", err.Error())
	}

	queue, err := ch.QueueDeclare(
		"message_queue", // queue name
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		ch.Close()
		conn.Close()
		log.Fatalf("Failed to declare a queue: %s", err.Error())
	}

	return RabbitMQ{
		Channel: ch,
		Queue:   queue,
	}
}

func (ra RabbitMQ) SendMessage(message models.QueueMessage) error {
	msgJSON, err := json.Marshal(message)
	if err != nil {
		return err
	}
	err = ra.Channel.PublishWithContext(context.Background(),
		"",
		ra.Queue.Name,
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        []byte(msgJSON),
		})
	if err != nil {
		log.Println("Failed to publish a message", err)
		return err
	}
	return nil
}
