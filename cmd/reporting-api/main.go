package main

import (
	"fmt"
	"net/http"
	"twitch_chat_analysis/thirdparty/redis"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	redis := redis.NewRedisInstance()

	r.GET("/message/list", func(c *gin.Context) {
		sender := c.Query("sender")
        receiver := c.Query("receiver")
		messages, err := redis.FetchFromRedis(fmt.Sprintf("%s:%s", sender, receiver))
		if err != nil{
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
            return
		}
		c.JSON(http.StatusOK, gin.H{"messages": messages})
	})

	r.Run(":8081")
}
