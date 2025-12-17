package main

import (
	"consumer/repository"
	"consumer/server"
	"consumer/services"
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()

	broker := os.Getenv("POSTS_BROKER")
	topic := os.Getenv("POSTS_TOPIC")
	groupID := os.Getenv("POSTS_GROUP_ID")

	repo := repository.KafkaConsumer{}
	err := repo.Init(broker, groupID, topic)
	if err != nil {
		fmt.Println("failed to initialize consumer repository")
		return
	}
	service := services.KafkaConsumerService{Consumer: repo}

	cserver := server.Init(&service)

	err = cserver.Run(os.Getenv("PORT"))
	if err != nil {
		fmt.Println("failed to initialize server")
	}

}
