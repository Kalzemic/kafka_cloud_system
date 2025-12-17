package services

import (
	"api/clients"
	"api/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

type ProducerService interface {
	CreatePost(c *gin.Context)
}

type APIProducerService struct {
	UClient clients.UserClient
	PClient clients.ProducerClient
}

func (service APIProducerService) CreatePost(c *gin.Context) {
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

	var post models.Post
	err = c.ShouldBindJSON(&post)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid input"})
		return
	}

	err = service.PClient.CreatePost(&post)
	if err != nil {
		switch err {
		case clients.ErrInvalidInput:
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid input"})
		default:
			c.JSON(http.StatusBadGateway, gin.H{"error": "failed to create post"})
		}
		return
	}
	c.JSON(http.StatusCreated, gin.H{"success": "post created succcessfully"})

}
