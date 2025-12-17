package services

import (
	"api/clients"
	"api/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

type ConsumerService interface {
	Poll(c *gin.Context)
}

type APIConsumerService struct {
	UClient clients.UserClient
	CClient clients.ConsumerClient
}

func (service APIConsumerService) Poll(c *gin.Context) {
	email := c.Param("email")
	password := c.Query("password")

	_, err := service.UClient.FindUser(email, password)
	if err != nil {
		switch err {
		case clients.ErrUserNotFound:
			c.JSON(http.StatusNotFound, gin.H{"error": "user doesnt exist"})
		default:
			c.JSON(http.StatusBadGateway, gin.H{"error": "failed to verify user"})
		}
		return
	}

	var req models.PollRequest
	err = c.ShouldBindJSON(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}
	var posts []models.Post
	posts, err = service.CClient.Poll(&req)
	if err != nil {
		switch err {
		case clients.ErrInvalidInput:
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid input"})
		default:
			c.JSON(http.StatusBadGateway, gin.H{"error": "failed to get posts"})
		}
		return
	}
	c.JSON(http.StatusOK, &posts)
}
