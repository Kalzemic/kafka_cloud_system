package repository

import (
	"consumer/models"
	"consumer/stream"
	"encoding/json"
	"fmt"
	"time"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

type Consumer interface {
	Init(broker, groupID, topic string) error

	Poll(max int, timeout time.Duration) ([]models.Post, error)
	Run(timeout time.Duration)
	Register() chan models.Post
	Unregister(chan models.Post)
	Close() error
}

type KafkaConsumer struct {
	Consumer *kafka.Consumer
	Topic    string
	GroupID  string
	Stream   *stream.Hub
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
	kc.Stream = stream.NewHub(1000)

	return nil
}

func (kc *KafkaConsumer) Run(timeout time.Duration) {

	go kc.Stream.Run()

	go func() {
		for {
			msg, err := kc.Consumer.ReadMessage(timeout)
			if err != nil {
				continue
			}
			var post models.Post
			if err = json.Unmarshal(msg.Value, &post); err != nil {
				fmt.Printf("Invalid message format: %v\n", err)
				continue
			}

			kc.Stream.Broadcast(post)
		}
	}()
}

func (kc *KafkaConsumer) Register() chan models.Post {
	return kc.Stream.Register(100)
}

func (kc *KafkaConsumer) Unregister(ch chan models.Post) {
	kc.Stream.Unregister(ch)
}

// func (kc *KafkaConsumer) Poll(max int, timeout time.Duration) ([]models.Post, error) {
// 	posts := make([]models.Post, 0, max)

// 	for len(posts) < max {
// 		msg, err := kc.Consumer.ReadMessage(timeout)
// 		if err != nil {
// 			break
// 		}
// 		var post models.Post
// 		if err = json.Unmarshal(msg.Value, &post); err != nil {
// 			fmt.Printf("Invalid message format: %v\n", err)
// 			continue
// 		}

// 		posts = append(posts, post)
// 	}
// 	return posts, nil
// }

func (kc *KafkaConsumer) Close() error {
	err := kc.Consumer.Close()
	if err != nil {
		return err
	}
	kc.Stream.Close()
	return nil
}
