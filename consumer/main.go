package main

import (
	"log"
	"time"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/protsack-stephan/wme-architecture-poc/pkg/schema"
)

const bootstrapServers = "localhost:29092"
const groupID = "main"

func main() {
	conn, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": bootstrapServers,
		"group.id":          groupID,
	})

	if err != nil {
		log.Panic(err)
	}

	defer conn.Close()

	if err := conn.SubscribeTopics([]string{schema.TopicPages, schema.TopicVersions}, nil); err != nil {
		log.Panic(err)
	}

	for {
		msg, err := conn.ReadMessage(time.Second * 5)

		if err != nil {
			log.Println(err)
			continue
		}

		log.Printf("key: '%s', val: '%s'", string(msg.Key), string(msg.Value))
	}
}
