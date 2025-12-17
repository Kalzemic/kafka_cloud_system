package main

import (
	"fmt"
	"kafka_service/repository"
	"kafka_service/server"
	"kafka_service/services"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()

	broker := os.Getenv("POSTS_BROKER")
	topic := os.Getenv("POSTS_TOPIC")

	repo := repository.KafkaProducer{}
	err := repo.Init(broker, topic)
	if err != nil {
		fmt.Println("failed to initialize kafka producer repository")
		return
	}
	service := services.KafkaProducerService{Producer: &repo}

	kserver := server.Init(&service)

	err = kserver.Run(os.Getenv("PORT"))
	if err != nil {
		fmt.Println("failed to run server")
		return
	}
}
