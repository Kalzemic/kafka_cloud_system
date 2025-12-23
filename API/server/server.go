package server

import (
	"api/services"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"time"
)

type APIServer struct {
	engine          *gin.Engine
	userService     services.UserService
	producerService services.ProducerService
	consumerService services.ConsumerService
}

func Init(us services.UserService, ps services.ProducerService, cs services.ConsumerService) *APIServer {
	serv := APIServer{engine: gin.Default(), userService: us, producerService: ps, consumerService: cs}

	serv.engine.Use(cors.New(cors.Config{
			AllowOrigins:     []string{"http://localhost:5173"},
			AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
			AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
			AllowCredentials: true,
			MaxAge:           12 * time.Hour,
		}))
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
		postsAPI.POST("produce/:email", serv.producerService.CreatePost)
		postsAPI.POST("poll/:email", serv.consumerService.Poll)
		postsAPI.GET("/listen/:email", serv.consumerService.Listen)
	}

	return &serv
}

func (serv *APIServer) Run(addr string) error {
	return serv.engine.Run(addr)
}
