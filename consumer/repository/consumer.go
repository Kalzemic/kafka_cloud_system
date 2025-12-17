package repository

import (
	"consumer/models"
	"encoding/json"
	"fmt"
	"time"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

type Consumer interface {
	Init(broker, groupID, topic string) error
	Poll(max int, timeout time.Duration) ([]models.Post, error)
	Listen()
	Close() error
}

type KafkaConsumer struct {
	Consumer *kafka.Consumer
	Topic    string
	GroupID  string
	Stream   chan models.Post
}

func (kc *KafkaConsumer) Init(broker, groupID, topic string) error {

	if broker == "" || groupID == "" || topic == "" {
		return fmt.Errorf("missing required env vars: BROKER, GROUP_ID, or TOPIC")
	}

	c, err := kafka.NewConsumer(&kafka.ConfigMap{"bootstrap.servers": broker, "group.id": groupID, "auto.offset.reset": "earliest"})
	if err != nil {
		return fmt.Errorf("create consumer: %w", err)
	}

	if err := c.SubscribeTopics([]string{topic}, nil); err != nil {
		return fmt.Errorf("subscribe to topic: %w", err)
	}

	kc.Consumer = c
	kc.Topic = topic
	kc.GroupID = groupID
	kc.Stream = make(chan models.Post, 1000)

	return nil
}

func (kc *KafkaConsumer) Listen() {

	go func() {

		fmt.Println("Listening Loop Initialized")

		for {
			msg, err := kc.Consumer.ReadMessage(-1)
			if err != nil {
				fmt.Printf("kafka read error%v\n", err)
				continue
			}

			var post models.Post
			if err := json.Unmarshal(msg.Value, &post); err != nil {
				fmt.Printf("Invalid message format%v\n", err)
				continue
			}

			select {
			case kc.Stream <- post:
				//
			default:
				<-kc.Stream
				kc.Stream <- post
			}

		}
	}()

}

func (kc *KafkaConsumer) Poll(max int, timeout time.Duration) ([]models.Post, error) {
	posts := make([]models.Post, 0, max)

	for len(posts) < max {
		msg, err := kc.Consumer.ReadMessage(timeout)
		if err != nil {
			break
		}
		var post models.Post
		if err = json.Unmarshal(msg.Value, &post); err != nil {
			fmt.Printf("Invalid message format: %v\n", err)
			continue
		}

		posts = append(posts, post)
	}
	return posts, nil
}

func (kc *KafkaConsumer) Close() error {
	err := kc.Consumer.Close()
	if err != nil {
		return err
	}
	close(kc.Stream)
	return nil
}
