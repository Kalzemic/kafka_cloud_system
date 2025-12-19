package services

import (
	"kafka_service/models"
	"kafka_service/repository"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type ProducerService interface {
	CreatePost(c *gin.Context)
}

type KafkaProducerService struct {
	Producer repository.Producer
}

func (service *KafkaProducerService) CreatePost(c *gin.Context) {
	var record models.Post

	err := c.ShouldBindJSON(&record)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid post format"})
		return
	}
	record.TimeStamp = time.Now()
	err = service.Producer.ProducePost(record)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to record post"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": "post recorded"})
}
