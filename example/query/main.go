package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

const url = "http://localhost:8088/query"
const query = `SELECT identifier, name, version__comment from pages_versions EMIT CHANGES;`

func main() {
	res, err := http.Post(url, "application/json", strings.NewReader(fmt.Sprintf(`{"ksql": "%s"}`, query)))

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
