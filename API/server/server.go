package server

import (
	"api/services"

	"github.com/gin-gonic/gin"
)

type APIServer struct {
	engine          *gin.Engine
	userService     services.UserService
	producerService services.ProducerService
	consumerService services.ConsumerService
}

func Init(us services.UserService, ps services.ProducerService, cs services.ConsumerService) *APIServer {
	serv := APIServer{engine: gin.Default(), userService: us, producerService: ps, consumerService: cs}

	userAPI := serv.engine.Group("/users")
	{
		userAPI.GET(":email", serv.userService.GetUserbyEmail)
		userAPI.POST("", serv.userService.CreateUser)
		userAPI.PUT(":email", serv.userService.UpdateUser)
		userAPI.GET("", serv.userService.GetAllUsers)
		userAPI.DELETE("", serv.userService.DeleteUsers)
	}

	postsAPI := serv.engine.Group("/posts")
	{
		postsAPI.POST("produce", serv.producerService.CreatePost)
		postsAPI.POST("poll", serv.consumerService.Poll)
	}

	return &serv
}

func (serv *APIServer) Run(addr string) error {
	return serv.engine.Run(addr)
}
