package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/protsack-stephan/dev-toolkit/lib/s3"
	"github.com/protsack-stephan/wme-architecture-poc/pkg/schema"
)

const bootstrapServers = "localhost:29092"
const groupID = "main"
const awsURL = "http://localhost:9000"
const awsID = "admin"
const awsKey = "password"
const awsRegion = "ap-northeast-1"
const awsBucket = "wme-data-bk"

func main() {
	store := s3.NewStorage(session.Must(session.NewSession(&aws.Config{
		Region:           aws.String(awsRegion),
		Credentials:      credentials.NewStaticCredentials(awsID, awsKey, ""),
		Endpoint:         aws.String(awsURL),
		DisableSSL:       aws.Bool(true),
		S3ForcePathStyle: aws.Bool(true),
	})), awsBucket)
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

		eMsg := new(schema.Message)

		if err := json.Unmarshal(msg.Value, eMsg); err != nil {
			log.Println(err)
			continue
		}

		path := fmt.Sprintf("%s.json", string(msg.Key))

		switch eMsg.Event.Type {
		case schema.EventTypeCreate, schema.EventTypeUpdate:
			if err := store.Put(path, bytes.NewReader(msg.Value)); err != nil {
				log.Printf("key: %s, err: %v\n", string(msg.Key), err)
			}
		case schema.EventTypeDelete:
			if err := store.Delete(path); err != nil {
				log.Printf("key: %s, err: %v\n", string(msg.Key), err)
			}
		default:
			log.Printf("key: %s, err: unknown event type '%s'\n", string(msg.Key), eMsg.Event.Type)
		}
	}
}
