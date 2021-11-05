package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

const url = "http://localhost:8088/ksql"

func main() {
	queries := []string{
		"CREATE STREAM pages(identifier INTEGER, name STRING, event STRUCT<uuid VARCHAR, type VARCHAR, date VARCHAR>, version STRUCT<identifier INTEGER>, date_modified VARCHAR, url VARCHAR, is_part_of STRUCT<identifier VARCHAR>, article_body ARRAY<STRUCT<text VARCHAR, encoding_format VARCHAR>>) WITH (kafka_topic='loc.pages.v1', value_format='json');",
		"CREATE STREAM versions(identifier INTEGER, name STRING, event STRUCT<uuid VARCHAR, type VARCHAR, date VARCHAR>, comment VARCHAR) WITH (kafka_topic='loc.versions.v1', value_format='json');",
		"CREATE STREAM pages_versions AS SELECT pages.name as name, pages.identifier as identifier, versions.identifier as version__identifier, versions.comment as version__comment FROM pages LEFT JOIN versions WITHIN 10 SECONDS ON pages.version->identifier = versions.identifier;",
		"CREATE TABLE pages_list (key STRING PRIMARY KEY, identifier INTEGER, name STRING, version STRUCT<identifier INTEGER>) WITH (kafka_topic='loc.pages.v1', value_format='json');",
		"CREATE TABLE queryable_pages_list AS SELECT * FROM pages_list;",
	}

	for _, q := range queries {
		res, err := http.Post(url, "application/json", strings.NewReader(fmt.Sprintf(`{"ksql": "%s"}`, q)))

		if err != nil {
			log.Panic(err)
		}

		data, err := ioutil.ReadAll(res.Body)
		res.Body.Close()

		if err != nil {
			log.Println(err)
		}

		if res.StatusCode != http.StatusOK {
			log.Println(string(data))
		}
	}
}
