package main

import (
	"context"
	"log"

	"github.com/protsack-stephan/wme-architecture-poc/ksqldb"
)

const url = "http://localhost:8088/"
const query = `SELECT identifier, name, version from pages EMIT CHANGES;`

func main() {
	cl := ksqldb.NewClient(url)

	err := cl.Push(context.Background(), &ksqldb.QueryRequest{SQL: query}, func(qr *ksqldb.QueryResponse, row ksqldb.Row) {
		log.Println(qr.QueryID, row)
	})

	if err != nil {
		log.Panic(err)
	}
}
