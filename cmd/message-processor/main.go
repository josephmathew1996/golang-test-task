package main

import (
	"encoding/json"
	"fmt"
	"log"
	"twitch_chat_analysis/pkg/models"
	"twitch_chat_analysis/thirdparty/rabbitmq"
	"twitch_chat_analysis/thirdparty/redis"
)

func main() {
	rabbitMQ := rabbitmq.NewRabbitMQInstance()
	redis := redis.NewRedisInstance()

	// this will continously consumes message from our queue
	for {
		msgs, err := rabbitMQ.Channel.Consume(
			rabbitMQ.Queue.Name,
			"",
			true,
			false,
			false,
			false,
			nil,
		)
		if err != nil {
			log.Fatalf("Failed to consume messages from: %s", err.Error())
		}

		//	TODO: this can be improve by implementing a worker-pool model and syncing with wait groups
		go func() {
			for d := range msgs {
				log.Printf("Received a message: %s", d.Body)
				var msg models.QueueMessage
				err := json.Unmarshal(d.Body, &msg)
				if err != nil {
					log.Println("Failed to parse the queue message: ", err)
					continue
				}
				err = redis.SaveToRedis(fmt.Sprintf("%s:%s", msg.Sender, msg.Receiver), d.Body)
				if err != nil {
					log.Println("Failed to send queue message to redis: ", err)
				}
				log.Println("Sent successfully to redis...")
			}
		}()
	}
}
