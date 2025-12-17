package services

import (
	"fmt"
	"net/http"
	"regexp"
	"strconv"
	"user_central/converter"
	"user_central/models"
	"user_central/storage"
	"user_central/validator"

	"github.com/gin-gonic/gin"
)

type UserService interface {
	CreateUser(c *gin.Context)
	UpdateUser(c *gin.Context)
	GetUserbyEmail(c *gin.Context)
	GetAllUsers(c *gin.Context)
	DeleteUsers(c *gin.Context)
}

type GinUserService struct {
	Repo storage.UserRepo
}

func (service *GinUserService) CreateUser(c *gin.Context) {
	var boundary models.UserBoundary
	err := c.ShouldBindJSON(&boundary)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	err = validator.ValidatePassword(boundary.Password)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	entity := converter.ConverttoEntity(boundary)
	err = service.Repo.CreateUser(&entity)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not save to database"})
		return
	}
	entityNew, err := service.Repo.FindUser(entity.Email, entity.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not load from database"})
		return
	}
	boundaryNew := converter.ConverttoBoundary(*entityNew)
	c.JSON(http.StatusCreated, &boundaryNew)

}

func (service *GinUserService) UpdateUser(c *gin.Context) {

	var boundary models.UserBoundary

	email := c.Param("email")
	password := c.Query("password")
	_, err := service.Repo.FindUser(email, password)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "no user found with matching credentials"})
		return
	}

	c.ShouldBindJSON(&boundary)
	entity := converter.ConverttoEntity(boundary)
	boundary.Email = email
	if err = service.Repo.UpdateUser(&entity); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not apply update"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": "user updated successfully"})

}

func (service *GinUserService) GetUserbyEmail(c *gin.Context) {
	email := c.Param("email")
	password := c.Query("password")
	boundary, err := service.Repo.FindUser(email, password)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "no user found with matching credentials"})
		return
	}
	boundary.Password = ""
	c.JSON(http.StatusOK, &boundary)
}

func (service *GinUserService) GetAllUsers(c *gin.Context) {
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

	var boundaries []models.UserBoundary
	var entities []models.UserEntity
	switch criteria {
	case "None":

		entities, err = service.Repo.GetAllUsers(page, size)

	case "byEmailDomain":
		domain := c.Query("value")
		if domain == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid parameters"})
			return
		}

		domainRegex := fmt.Sprintf("@%s$", regexp.QuoteMeta(domain))
		entities, err = service.Repo.GetbyEmailDomain(domainRegex, page, size)

	case "byRole":
		role := c.Query("value")
		if role == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid parameters"})
			return
		}
		entities, err = service.Repo.GetUsersbyRoles(role, page, size)

	case "byRegistrationToday":
		entities, err = service.Repo.GetUsersbyRegistrationToday(page, size)
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to retrieve users"})
		return
	}

	for _, entity := range entities {
		boundary := converter.ConverttoBoundary(entity)
		boundary.Password = ""
		boundaries = append(boundaries, boundary)
	}
	c.JSON(http.StatusOK, boundaries)
}

func (service *GinUserService) DeleteUsers(c *gin.Context) {
	err := service.Repo.DeleteUsers()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete users"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": "users deleted successfully"})
}
