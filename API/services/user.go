package services

import (
	"api/clients"
	"api/models"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type UserService interface {
	CreateUser(c *gin.Context)
	UpdateUser(c *gin.Context)
	GetUserbyEmail(c *gin.Context)
	GetAllUsers(c *gin.Context)
	DeleteUsers(c *gin.Context)
}

type APIUserService struct {
	client clients.UserClient
}

func (us *APIUserService) CreateUser(c *gin.Context) {
	var req models.UserRequest
	err := c.ShouldBindJSON(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request credentials"})
		return
	}
	resp, err := us.client.CreateUser(&req)
	if err != nil {
		switch err {
		case clients.ErrInvalidInput:
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid input"})
		case clients.ErrUserAlreadyExists:
			c.JSON(http.StatusConflict, gin.H{"error": "user already exists"})
		default:
			c.JSON(http.StatusBadGateway, gin.H{"error": "failed to create user"})

		}
		return
	}
	c.JSON(http.StatusCreated, resp)
}

func (us *APIUserService) UpdateUser(c *gin.Context) {
	var req models.UserRequest
	err := c.ShouldBindJSON(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request credentials"})
		return
	}

	email := c.Param("email")
	req.Email = email

	err = us.client.UpdateUser(&req)
	if err != nil {
		switch err {
		case clients.ErrInvalidInput:
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid input"})
		case clients.ErrUserNotFound:
			c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		default:
			c.JSON(http.StatusBadGateway, gin.H{"error": "failed to update user"})
		}
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": "user updated successfully"})
}

func (us *APIUserService) GetUserbyEmail(c *gin.Context) {
	email := c.Param("email")
	password := c.Query("password")
	boundary, err := us.client.FindUser(email, password)
	if err != nil {
		switch err {

		case clients.ErrUserNotFound:
			c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})

		default:
			c.JSON(http.StatusBadGateway, gin.H{"error": "failed to fetch user"})
		}
		return
	}
	c.JSON(http.StatusOK, &boundary)
}

func (us *APIUserService) GetAllUsers(c *gin.Context) {
	criteria := c.DefaultQuery("criteria", "None")

	pageStr := c.DefaultQuery("page", "1")
	page, err := strconv.Atoi(pageStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid parameters"})
		return
	}

	sizeStr := c.DefaultQuery("size", "10")
	size, err := strconv.Atoi(sizeStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid parameters"})
		return
	}

	var responses []models.UserResponse
	switch criteria {
	case "None":

		responses, err = us.client.GetAllUsers(page, size)

	case "byEmailDomain":
		domain := c.Query("value")
		if domain == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid parameters"})
			return
		}

		responses, err = us.client.GetbyEmailDomain(domain, page, size)

	case "byRole":
		role := c.Query("value")
		if role == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid parameters"})
			return
		}
		responses, err = us.client.GetUsersbyRoles(role, page, size)

	case "byRegistrationToday":
		responses, err = us.client.GetUsersbyRegistrationToday(page, size)

	default:
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid criteria"})
		return
	}

	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": "failed to retrieve users"})
		return
	}

	c.JSON(http.StatusOK, responses)
}

func (us *APIUserService) DeleteUsers(c *gin.Context) {
	err := us.client.DeleteUsers()
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": "failed to delete users"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": "users deleted successfully"})
}
