package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/google/uuid"
	"github.com/protsack-stephan/mediawiki-api-client"
	eventstream "github.com/protsack-stephan/mediawiki-eventstream-client"
	"github.com/protsack-stephan/wme-architecture-poc/pkg/schema"
)

const database = "enwiki"
const url = "https://en.wikipedia.org"
const bootstrapServers = "localhost:29092"

func main() {
	ctx := context.Background()
	since := time.Now()
	streams := eventstream.NewClient()
	mediawiki := mediawiki.NewClient(url)
	producer, err := kafka.NewProducer(&kafka.ConfigMap{
		"bootstrap.servers": bootstrapServers,
	})

	if err != nil {
		log.Panic(err)
	}

	defer producer.Close()

	stream := streams.RevisionCreate(ctx, since, func(evt *eventstream.RevisionCreate) {
		if evt.Data.Database == database {
			data, err := mediawiki.PagesData(ctx, evt.Data.PageTitle)

			if err != nil {
				log.Println(err)
				return
			}

			pageData, ok := data[evt.Data.PageTitle]

			if !ok {
				return
			}

			html, err := mediawiki.PageHTML(ctx, evt.Data.PageTitle)

			if err != nil {
				log.Println(err)
				return
			}

			page := new(schema.Page)
			page.Name = pageData.Title
			page.Identifier = pageData.PageID
			page.DateModified = pageData.Touched
			page.URL = fmt.Sprintf("%s/wiki/%s", url, evt.Data.PageTitle)
			page.IsPartOf = &schema.Project{
				Identifier: database,
			}
			page.Version = &schema.Version{
				Identifier: pageData.Revisions[0].RevID,
			}

			page.ArticleBody = append(page.ArticleBody, &schema.ArticleBody{
				Text:           string(html),
				EncodingFormat: "text/html",
			})

			page.ArticleBody = append(page.ArticleBody, &schema.ArticleBody{
				Text:           pageData.Revisions[0].Slots.Main.Content,
				EncodingFormat: pageData.Revisions[0].Slots.Main.Contentformat,
			})

			page.Event = new(schema.Event)
			page.Event.UUID = uuid.NewString()
			page.Event.Date = time.Now().UTC()
			page.Event.Type = schema.EventTypeUpdate

			pmsg := kafka.Message{
				TopicPartition: kafka.TopicPartition{Topic: &schema.TopicPages, Partition: 0},
				Key:            []byte(fmt.Sprintf("pages/%s/%s", evt.Data.Database, evt.Data.PageTitle)),
			}

			if pmsg.Value, err = json.Marshal(page); err != nil {
				log.Println(err)
				return
			}

			if err := producer.Produce(&pmsg, nil); err != nil {
				log.Println(err)
				return
			}

			version := new(schema.Version)
			version.Comment = pageData.Revisions[0].Comment
			version.Identifier = pageData.Revisions[0].RevID

			version.Event = new(schema.Event)
			version.Event.UUID = uuid.NewString()
			version.Event.Date = time.Now().UTC()
			version.Event.Type = schema.EventTypeCreate

			vmsg := kafka.Message{
				TopicPartition: kafka.TopicPartition{Topic: &schema.TopicVersions, Partition: 0},
				Key:            []byte(fmt.Sprintf("versions/%s/%d", evt.Data.Database, evt.Data.RevID)),
			}

			if vmsg.Value, err = json.Marshal(version); err != nil {
				log.Println(err)
				return
			}

			if err := producer.Produce(&vmsg, nil); err != nil {
				log.Println(err)
				return
			}
		}
	})

	for err := range stream.Sub() {
		if err != nil {
			log.Println(err)
		}
	}
}
