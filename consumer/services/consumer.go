package services

import (
	"consumer/models"
	"consumer/repository"
	"net/http"

	"github.com/gin-gonic/gin"
)

type ConsumerService interface {
	Listen(c *gin.Context)
	Poll(c *gin.Context)
}

type KafkaConsumerService struct {
	Consumer repository.KafkaConsumer
}

func (service KafkaConsumerService) Listen(c *gin.Context) {
	panic("unimplemented")
}

func (service *KafkaConsumerService) Poll(c *gin.Context) {
	var req models.PollRequest

	err := c.ShouldBindJSON(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid poll request"})
	}

	posts := make([]models.Post, req.MaxPosts)

	posts, err = service.Consumer.Poll(req.MaxPosts, req.MaxDuration)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusFound, posts)
}
