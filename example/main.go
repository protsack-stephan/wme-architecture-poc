package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

func main() {
	res, err := http.Post("http://localhost:8088/query", "application/json", strings.NewReader(`{"ksql": "SELECT identifier, name, event, version->identifier as versionID FROM pages EMIT CHANGES;"}`))
	// res, err := http.Post("http://localhost:8088/query", "application/json", strings.NewReader(`{"ksql": "SELECT identifier, comment, event FROM versions EMIT CHANGES;"}`))

	if err != nil {
		log.Panic(err)
	}

	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		data, err := ioutil.ReadAll(res.Body)

		if err != nil {
			log.Panic(err)
		}

		log.Panic(string(data))
	}

	scn := bufio.NewScanner(res.Body)

	for scn.Scan() {
		if len(scn.Text()) > 0 {
			fmt.Println(scn.Text())
		}
	}
}
