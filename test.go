package main

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/segmentio/kafka-go"
)

func main2() {
	// to produce messages
	topic := "tpc-req-master"
	partition := 0

	conn, err := kafka.DialLeader(context.Background(), "tcp", "localhost:9092", topic, partition)
	if err != nil {
		log.Fatal("failed to dial leader:", err)
	}

	data := make(map[string]interface{})
	data["companyId"] = "1001"
	data["companyName"] = "PT Salma 1"
	respData, err := json.Marshal(data)

	// s := fmt.Sprintf("%v", data)
	// req["data"] = data

	conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
	_, err = conn.WriteMessages(
		kafka.Message{
			Value: []byte(respData),
		},
	)
	if err != nil {
		log.Fatal("failed to write messages:", err)
	}

	if err := conn.Close(); err != nil {
		log.Fatal("failed to close writer:", err)
	}
}
