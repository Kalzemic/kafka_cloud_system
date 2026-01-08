package server

import (
	"consumer/services"
	"time"

	"github.com/gin-gonic/gin"
)

type ConsumerServer struct {
	Engine   *gin.Engine
	Consumer services.ConsumerService
}

func Init(consumer services.ConsumerService) ConsumerServer {
	server := ConsumerServer{gin.Default(), consumer}

	api := server.Engine.Group("posts")
	{
		// api.POST("poll", server.Consumer.Poll)
		api.GET("listen", server.Consumer.Listen)
	}

	return server
}

func (cserver ConsumerServer) Run(addr string) error {
	cserver.Consumer.Run(time.Second * 2)
	return cserver.Engine.Run(addr)
}
