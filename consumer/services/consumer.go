package services

import (
	
	"consumer/repository"
	"io"
	"time"

	"github.com/gin-gonic/gin"
)

type ConsumerService interface {
	Run(timeout time.Duration)
	Listen(c *gin.Context)
	// Poll(c *gin.Context)
}

type KafkaConsumerService struct {
	Consumer *repository.KafkaConsumer
}

func (service *KafkaConsumerService) Run(timeout time.Duration) {
	service.Consumer.Run(timeout)
}

func (service KafkaConsumerService) Listen(c *gin.Context) {
	ch := service.Consumer.Register()
	defer service.Consumer.Unregister(ch)
	c.Header("Content-Type", "text/event-stream")
	c.Header("Cache-Control", "no-cache")
	c.Header("Connection", "keep-alive")
	c.Stream(func(w io.Writer) bool {
		select {
		case post, ok := <-ch:
			if !ok {
				return false
			}
			c.SSEvent("post", post)
			return true
		case <-c.Request.Context().Done():
			return false
		}

	})
}

// func (service *KafkaConsumerService) Poll(c *gin.Context) {
// 	var req models.PollRequest

// 	err := c.ShouldBindJSON(&req)
// 	if err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid poll request"})
// 		return
// 	}

// 	var posts []models.Post

// 	posts, err = service.Consumer.Poll(req.MaxPosts, req.MaxDuration)
// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
// 		return
// 	}
// 	c.JSON(http.StatusOK, posts)
// }
