package redis

import (
	"encoding/json"
	"log"
	"twitch_chat_analysis/pkg/models"

	"github.com/go-redis/redis"
)

type RedisInstance struct {
	client *redis.Client
}

func NewRedisInstance() RedisInstance {
	client := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})
	_, err := client.Ping().Result()
	if err != nil {
		log.Fatalf("Failed to connect to Redis: %s", err)
	}

	log.Println("Connected to Redis")
	return RedisInstance{
		client: client,
	}
}

func (rc RedisInstance) SaveToRedis(key string, message []byte) error {
	err := rc.client.LPush(key, message).Err()
	if err != nil {
		return err
	}
	return nil
}

func (rc RedisInstance) FetchFromRedis(key string) ([]models.QueueMessage, error) {
	// We are using LRange because we used LPush in our mesage processor service to push the message to redis
	// So the elements are pushed to the beginning or the left side of a list.
	// That means last message will be present in the leftmost position.
	// So to receive the message in chronological descending order we can use LRange from index 0 till the end of the list
	msgs, err := rc.client.LRange(key, 0, -1).Result()
	if err != nil {
		return nil, err
	}
	var messages []models.QueueMessage
	for _, msgJSON := range msgs {
		var msg models.QueueMessage
		if err := json.Unmarshal([]byte(msgJSON), &msg); err != nil {
			return nil, err
		}
		messages = append(messages, msg)
	}
	return messages, err
}
