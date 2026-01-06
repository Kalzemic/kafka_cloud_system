package services

import (
	"api/clients"
	"api/models"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
)

type ConsumerService interface {
	Poll(c *gin.Context)
	Listen(c *gin.Context)
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

func (service APIConsumerService) Listen(c *gin.Context) {
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
	ctx := c.Request.Context() 

	ch, err := service.CClient.Listen(ctx)
	if err != nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{"error": "stream unavailable"})
		return
	}

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

		case <-ctx.Done():

			return false
		}
	})
}
