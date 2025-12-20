package main

import (
	"api/clients"
	"api/server"
	"api/services"
	"fmt"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

func main() {

	godotenv.Load()

	userURL := os.Getenv("USER_URL")
	producerURL := os.Getenv("PRODUCER_URL")
	consumerURL := os.Getenv("CONSUMER_URL")

	uc := &clients.APIUserClient{Client: &http.Client{Timeout: 0}, BaseURL: userURL}
	pc := &clients.APIProducerClient{Client: &http.Client{Timeout: 0}, BaseURL: producerURL}
	cc := &clients.APIConsumerClient{Client: &http.Client{Timeout: 0}, BaseURL: consumerURL}

	us := &services.APIUserService{Client: uc}
	ps := &services.APIProducerService{UClient: uc, PClient: pc}
	cs := &services.APIConsumerService{UClient: uc, CClient: cc}

	serv := server.Init(us, ps, cs)

	err := serv.Run(os.Getenv("PORT"))
	if err != nil {
		fmt.Printf("failed to run server %s\n", err.Error())
	}
}
