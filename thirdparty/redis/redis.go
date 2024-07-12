package redis

import (
	"log"

	"github.com/go-redis/redis"
)


type RedisInstance struct {
	client  *redis.Client
}

func NewRedisInstance() RedisInstance {
    client := redis.NewClient(&redis.Options{
        Addr:     "localhost:6379",
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