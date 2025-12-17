package server

import (
	"kafka_service/services"

	"github.com/gin-gonic/gin"
)

type ProducerServer struct {
	Engine   *gin.Engine
	Producer services.ProducerService
}

func Init(producer services.ProducerService) ProducerServer {
	server := ProducerServer{gin.Default(), producer}

	api := server.Engine.Group("posts")
	{
		api.POST("produce", server.Producer.CreatePost)
	}

	return server
}

func (kserver ProducerServer) Run(addr string) error {
	return kserver.Engine.Run(addr)
}
