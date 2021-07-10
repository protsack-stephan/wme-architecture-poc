package main

import (
	"context"
	"log"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/protsack-stephan/wme-architecture-poc/pkg/schema"
)

const bootstrapServers = "localhost:29092"

func main() {
	ctx := context.Background()
	adm, err := kafka.NewAdminClient(&kafka.ConfigMap{
		"bootstrap.servers": bootstrapServers,
	})

	if err != nil {
		log.Panic(err)
	}

	defer adm.Close()

	config := map[string]string{
		"cleanup.policy":            "compact,delete",
		"min.cleanable.dirty.ratio": "0.01",
		"delete.retention.ms":       "100",
		"segment.ms":                "100",
		"max.message.bytes":         "20971520",
	}

	params := []kafka.TopicSpecification{
		{
			Topic:             schema.TopicPages,
			NumPartitions:     1,
			ReplicationFactor: 1,
			Config:            config,
		},
		{
			Topic:             schema.TopicVersions,
			NumPartitions:     1,
			ReplicationFactor: 1,
			Config:            config,
		},
	}

	statuses, err := adm.CreateTopics(ctx, params)

	if err != nil {
		log.Panic(err)
	}

	for _, status := range statuses {
		log.Printf("topic: %s, err: %s\n", status.Topic, status.Error)
	}
}
