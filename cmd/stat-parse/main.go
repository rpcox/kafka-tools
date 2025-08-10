package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
)

func main() {
	fh, err := os.Open("data/consumer.stats")
	if err != nil {
		log.Fatal(err)
	}

	data, err := io.ReadAll(fh)
	if err != nil {
		log.Fatal(err)
	}

	var stats LibRdKafkaStats
	err = json.Unmarshal(data, &stats)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(stats)
}
