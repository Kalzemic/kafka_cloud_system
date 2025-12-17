package server

import (
	"user_central/services"

	"github.com/gin-gonic/gin"
)

type UserServer struct {
	engine  *gin.Engine
	Service services.UserService
}

func ServerInit(service services.UserService) *UserServer {
	server := &UserServer{engine: gin.Default(), Service: service}

	usersAPI := server.engine.Group("/users")
	{
		usersAPI.POST("", service.CreateUser)
		usersAPI.PUT(":email", service.UpdateUser)
		usersAPI.GET(":email", service.GetUserbyEmail)
		usersAPI.GET("", service.GetAllUsers)
		usersAPI.DELETE("", service.DeleteUsers)
	}

	postsAPI := server.engine.Group("/posts")
	{
		postsAPI.POST("", func(c *gin.Context) {})
	}

	return server
}

func (server *UserServer) Run(addr string) error {
	return server.engine.Run(addr)
}
