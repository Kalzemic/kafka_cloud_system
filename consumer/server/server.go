package server

import (
	"consumer/services"

	"github.com/gin-gonic/gin"
)

type ConsumerServer struct {
	Engine   *gin.Engine
	Producer services.ConsumerService
}

func Init(producer services.ConsumerService) ConsumerServer {
	server := ConsumerServer{gin.Default(), producer}

	api := server.Engine.Group("posts")
	{
		api.POST("poll", server.Producer.Poll)
	}

	return server
}

func (cserver ConsumerServer) Run(addr string) error {
	return cserver.Engine.Run(addr)
}
