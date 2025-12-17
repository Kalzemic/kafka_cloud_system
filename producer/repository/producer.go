package repository

import (
	"encoding/json"
	"fmt"
	"kafka_service/models"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

type Producer interface {
	Init(broker, topic string) error
	ProducePost(post models.Post) error
	Close()
}

type KafkaProducer struct {
	Producer *kafka.Producer
	Topic    string
	//DeliveryChan chan kafka.Event
}

func (kp *KafkaProducer) Init(broker, topic string) error {

	if broker == "" || topic == "" {
		return fmt.Errorf("missing required env vars: BROKER or TOPIC")
	}

	prod, err := kafka.NewProducer(&kafka.ConfigMap{
		"bootstrap.servers": broker,
	})

	if err != nil {
		return err
	}
	kp.Producer = prod
	kp.Topic = topic
	//kp.DeliveryChan = make(chan kafka.Event, 10000)

	fmt.Printf("Kafka producer initialized (broker=%s, topic=%s)\n", broker, topic)
	return nil
}

func (kp *KafkaProducer) ProducePost(post models.Post) error {
	data, err := json.Marshal(post)
	if err != nil {
		return fmt.Errorf("failed to marshal post: %w", err)
	}

	if err = kp.Producer.Produce(&kafka.Message{TopicPartition: kafka.TopicPartition{Partition: kafka.PartitionAny, Topic: &kp.Topic}, Value: data}, nil); err != nil {
		return fmt.Errorf("produce: %w", err)
	}
	return nil

}

func (kp *KafkaProducer) Close() {
	fmt.Println("Closing Kafka producer")
	kp.Producer.Flush(15 * 1000)
	kp.Producer.Close()
	//close(kp.DeliveryChan)
}

// go func() {
// 	for e := range kp.DeliveryChan {
// 		switch ev := e.(type) {
// 		case *kafka.Message:
// 			if ev.TopicPartition.Error != nil {
// 				fmt.Printf("Delivery Failed %s\n", ev.TopicPartition.Error.Error())
// 			} else {
// 				fmt.Printf(" Delivered %s [%d] @ %v\n",
// 					*ev.TopicPartition.Topic, ev.TopicPartition.Partition, ev.TopicPartition.Offset)
// 			}
// 		default:
// 			fmt.Printf("Kafka event:%v\n ", ev)
// 		}
// 	}
// 	fmt.Println("Delivery channel listener exiting.")
// }()
