package main

import (
	"net/http"
	"twitch_chat_analysis/pkg/models"
	"twitch_chat_analysis/thirdparty/rabbitmq"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	rabbitMQ := rabbitmq.NewRabbitMQInstance()

	r.GET("/test", func(c *gin.Context) {
		c.JSON(200, "worked")
	})

	r.POST("/message", func(c *gin.Context) {
		var msg models.QueueMessage
        if err := c.ShouldBindJSON(&msg); err != nil {
            c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
            return
        }

        err := rabbitMQ.SendMessage(msg)
        if err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to push message to queue"})
            return
        }

        c.JSON(http.StatusOK, gin.H{"status": "Message pushed to queued successfully"})
	})

	r.Run()
}
