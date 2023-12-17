package main

import (
	"WBL0/internal/model"
	"encoding/json"
	"github.com/nats-io/stan.go"
	"io"
	"log"
	"os"
)

const (
	producerID  = "producer"
	stanChannel = "channel"
	clusterID   = "test-cluster"
)

func main() {

	sc, err := stan.Connect(clusterID, producerID)
	if err != nil {
		log.Fatal(err)
	}
	defer sc.Close()

	jsonTestData, err := os.Open("test_data.json")
	if err != nil {
		log.Printf("error while openning file %v", err)
	}
	defer jsonTestData.Close()

	byteValue, _ := io.ReadAll(jsonTestData)

	var testOrders []model.ValidOrderTemplate

	err = json.Unmarshal(byteValue, &testOrders)
	if err != nil {
		log.Printf("error while parsing file %v", err)
	}

	for _, order := range testOrders {
		jsonOrder, _ := json.Marshal(order)

		err = sc.Publish(stanChannel, jsonOrder)
		if err != nil {
			log.Println(err)
		}
	}

	err = sc.Publish(stanChannel, []byte("Invalid JSON"))
	if err != nil {
		log.Println(err)
	}

	log.Println("all messages sent")
}
