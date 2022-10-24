package main

import (
	"context"
	"encoding/json"
	"kafka-user-dc/models"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	uuid "github.com/satori/go.uuid"

	conf "kafka-user-dc/config"
	con "kafka-user-dc/controllers"
)

type Respon struct {
	Expired float64     `json:"expired"`
	Data    interface{} `json:"data"`
	Uuid    uuid.UUID
}

type MessageKafka struct {
	Topic string `json:"topic"`
	Proc  string `json:"proc"`
	Data  string `json:"data"`
	From  string `json:"from"`
	Exp   int64  `json:"exp"`
}

type MessageWebsocket struct {
	Type  int    `json:"type"`
	Topic string `json:"topic"`
	Proc  string `json:"proc"`
	Data  string `json:"data"`
}

var UuidReq uuid.UUID

func main() {

	consumenUser()
}

func consumenUser() {
	log.Printf("consumenUser")

	// createTopic("user-request", conf.GetBroker())
	log.Printf("brokerServer : %v", conf.GetBroker())
	//create consumer
	var topics []string
	topics = append(topics, "user-request")
	sigchan := make(chan os.Signal, 1)
	signal.Notify(sigchan, syscall.SIGINT, syscall.SIGTERM)
	//earliest: automatically reset the offset to the earliest offset
	//latest: automatically reset the offset to the latest offset
	//none: throw exception to the consumer if no previous offset is found for the consumer's group
	c, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers":               conf.GetBroker(),
		"group.id":                        "user",
		"session.timeout.ms":              6000,
		"go.events.channel.enable":        true,
		"go.application.rebalance.enable": true,
		"auto.offset.reset":               "latest",
	})
	if err != nil {
		log.Printf("Failed to create consumer: %s\n", err)
		os.Exit(1)
	}
	log.Printf("Created Consumer %v\n", c)
	err = c.SubscribeTopics(topics, nil)
	run := true
	for run {
		select {
		case sig := <-sigchan:
			log.Printf("Caught signal %v: terminating\n", sig)
			run = false
		case ev := <-c.Events():
			switch e := ev.(type) {
			case kafka.AssignedPartitions:
				log.Printf("%% %v\n", e)
				c.Assign(e.Partitions)
			case kafka.RevokedPartitions:
				log.Printf("%% %v\n", e)
				c.Unassign()
			case *kafka.Message:
				// log.Printf("%% Message on %s:\n%s\n", e.TopicPartition, string(e.Value))
				s := string(e.Value)
				log.Printf("message kafka : %v", s)
				// data := make(map[string]interface{})
				data := &MessageKafka{}
				json.Unmarshal([]byte(s), &data)

				// sRequest, _ := json.Marshal(data.Data)
				sRequest := ""
				if data.Proc == "new" {

					userData := &models.User{}
					json.Unmarshal([]byte(data.Data), &userData)
					log.Printf("userData : %v", userData.Username)
					outputData, _ := models.CreateUser(userData)
					byteRespon, err := json.Marshal(outputData)
					if err != nil {
						log.Printf("Error : %s\n", err)
						// return docsError, err
						sRequest = "error"
					} else {
						sRequest = string(byteRespon)
					}

				} else if data.Proc == "edit" {
					userData := &models.User{}
					json.Unmarshal([]byte(data.Data), &userData)
					log.Printf("userData : %v", userData.Username)
					outputData, _ := models.UpdateUser(userData)
					byteRespon, err := json.Marshal(outputData)
					if err != nil {
						log.Printf("Error : %s\n", err)
						// return docsError, err
						sRequest = "error"
					} else {
						sRequest = string(byteRespon)
					}

				} else if data.Proc == "delete" {
					userData := &models.User{}
					json.Unmarshal([]byte(data.Data), &userData)
					log.Printf("userData : %v", userData.Username)

					outputData, _ := models.DeleteUser(userData)
					byteRespon, err := json.Marshal(outputData)
					if err != nil {
						log.Printf("Error : %s\n", err)
						// return docsError, err
						sRequest = "error"
					} else {
						sRequest = string(byteRespon)
					}

				} else if data.Proc == "search" {
					userData := &models.User{}
					json.Unmarshal([]byte(data.Data), &userData)
					log.Printf("userData : %v", userData.Username)
					outputData, _ := models.SearchUser(userData)
					byteRespon, err := json.Marshal(outputData)
					if err != nil {
						log.Printf("Error : %s\n", err)
						// return docsError, err
						sRequest = "error"
					} else {
						sRequest = string(byteRespon)
					}
				}

				log.Printf("sRequest : %v", sRequest)
				data.Topic = "user"
				data.Data = sRequest
				log.Printf("LANJUT PROSES UNTUK KE DATA")
				err = sendToBroker("user-respon", data, "test", "test")
				if err != nil {
					// handle this errors
					log.Printf("Error sendToBroker %v\n", err)
				}

			case kafka.PartitionEOF:
				log.Printf("%% Reached %v\n", e)
			case kafka.Error:
				log.Printf("Error: %v\n", e)
				run = false

			}
			log.Printf("cek...\n")
		}
		log.Printf("cek-2...\n")
	}

	log.Printf("Closing consumer\n")
	log.Printf("End at %v\n", time.Now())
	c.Close()
	os.Exit(3)

}

func sendToBroker(topic string, data *MessageKafka, key string, valHeaderBroker string) error {

	procdr, err := kafka.NewProducer(&kafka.ConfigMap{
		"bootstrap.servers":   conf.GetBroker(),
		"delivery.timeout.ms": 1000,
	})
	if err != nil {
		log.Printf("Failed to create producer: %s\n", err)
		os.Exit(1)
	}
	log.Printf("Created Producer %v\n", procdr)

	// timeProcess := time.Now()

	if procdr != nil {
		dataJson, err := json.Marshal(data)
		if err != nil {
			log.Printf("Error parse to JSON : %v\n", err)
		}

		//producer channel
		doneChan := make(chan bool)
		go func() {
			defer close(doneChan)
			for e := range procdr.Events() {
				switch ev := e.(type) {
				case *kafka.Message:
					m := ev
					if m.TopicPartition.Error != nil {
						log.Printf("Sent failed: %v\n", m.TopicPartition.Error)
					} else {
						log.Printf("Sent message to topic %s [%d] at offset %v\n",
							*m.TopicPartition.Topic, m.TopicPartition.Partition, m.TopicPartition.Offset)
					}
					return

				default:
					log.Printf("Ignored event: %s\n", ev)
				}
			}
		}()
		value := string(dataJson)
		//log.Printf("dataJson: %v\n", value)
		procdr.ProduceChannel() <- &kafka.Message{
			TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
			Value:          []byte(value),
			Headers:        []kafka.Header{{Key: key, Value: []byte(valHeaderBroker)}},
		}
		// wait for delivery report goroutine to finish
		_ = <-doneChan

	} else {
		//if database not defined, display packet in log
		log.Printf("Message Broker not defined\n")
	}

	// duration := time.Since(timeProcess)
	log.Printf("Sent To Broker\n")

	return nil
}

func message(code uint, status bool, message string) map[string]interface{} {
	if strings.TrimSpace(strings.ToLower(message)) == "success" {
		return map[string]interface{}{"code": code, "status": status, "message": http.StatusText(http.StatusOK)}
	} else {
		return map[string]interface{}{"code": code, "status": status, "message": message}
	}
}

func responBroker(topic string, broker string, group string) {

	//create consumer
	var topics []string
	topics = append(topics, topic)
	sigchan := make(chan os.Signal, 1)
	signal.Notify(sigchan, syscall.SIGINT, syscall.SIGTERM)
	//earliest: automatically reset the offset to the earliest offset
	//latest: automatically reset the offset to the latest offset
	//none: throw exception to the consumer if no previous offset is found for the consumer's group
	c, err := kafka.NewConsumer(&kafka.ConfigMap{
		// "bootstrap.servers":  broker,
		// "group.id":           group,
		// "session.timeout.ms": 6000,
		// // "go.events.channel.enable":        true,
		// // "go.application.rebalance.enable": true,
		// "default.topic.config": kafka.ConfigMap{"auto.offset.reset": "latest"},
		// // "auto.offset.reset":        "earliest",
		// // "enable.auto.offset.store": false,

		"bootstrap.servers":               broker,
		"group.id":                        group,
		"session.timeout.ms":              6000,
		"go.events.channel.enable":        true,
		"go.application.rebalance.enable": true,
		"default.topic.config":            kafka.ConfigMap{"auto.offset.reset": "earliest"},
	})
	if err != nil {
		log.Printf("Failed to create consumer: %s\n", err)
		os.Exit(1)
	}
	log.Printf("Created Consumer %v\n", c)
	err = c.SubscribeTopics(topics, nil)
	run := true
	for run == true {
		select {
		case sig := <-sigchan:
			log.Printf("Caught signal %v: terminating\n", sig)
			run = false
		case ev := <-c.Events():
			// ev = c.Poll(10)
			switch e := ev.(type) {
			case kafka.AssignedPartitions:
				log.Printf("%% %v\n", e)
				c.Assign(e.Partitions)
			case kafka.RevokedPartitions:
				log.Printf("%% %v\n", e)
				c.Unassign()
			case *kafka.Message:
				log.Printf("%% Message on %s:\n%s\n", e.TopicPartition, string(e.Value))
				s := string(e.Value)
				// data := make(map[string]interface{})
				// json.Unmarshal([]byte(s), &data)

				data := &con.User{}
				json.Unmarshal([]byte(s), data)

				con.CreateUser(data)
			case kafka.PartitionEOF:
				log.Printf("%% Reached %v\n", e)
			case kafka.Error:
				log.Printf("%% NewConsumer-Error: %v\n", e)
				run = false
			// case kafka.ErrorCode:
			// 	log.Printf("kafka.ErrorCode: %v\n", e)

			default:
				// log.Printf("keluar gak nih ?\n")
				run = false
			}
		}
	}

	log.Printf("Closing consumer\n")
	log.Printf("End at %v\n", time.Now())
	c.Close()

}

func createTopic(topic string, bootstrapServers string) {

	adminClient, err := kafka.NewAdminClient(&kafka.ConfigMap{
		"bootstrap.servers": bootstrapServers,
		// "broker.version.fallback": "0.10.0.0",
		// "api.version.fallback.ms": 0,
		// "sasl.mechanisms": "PLAIN",
		// "security.protocol":       "SASL_SSL",
		// "sasl.username":           ccloudAPIKey,
		// "sasl.password":           ccloudAPISecret
	})

	if err != nil {
		log.Printf("Failed to create Admin client: %s\n", err)
		os.Exit(1)
	}

	// Contexts are used to abort or limit the amount of time
	// the Admin call blocks waiting for a result.
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Create topics on cluster.
	// Set Admin options to wait for the operation to finish (or at most 60s)
	maxDuration, err := time.ParseDuration("60s")
	if err != nil {
		panic("time.ParseDuration(60s)")
	}

	results, err := adminClient.CreateTopics(ctx,
		[]kafka.TopicSpecification{{
			Topic:             topic,
			NumPartitions:     1,
			ReplicationFactor: 3}},
		kafka.SetAdminOperationTimeout(maxDuration))

	if err != nil {
		log.Printf("Problem during the topic creation: %v\n", err)
		os.Exit(1)
	}

	// Check for specific topic errors
	for _, result := range results {
		if result.Error.Code() != kafka.ErrNoError &&
			result.Error.Code() != kafka.ErrTopicAlreadyExists {
			log.Printf("Topic creation failed for %s: %v",
				result.Topic, result.Error.String())
			os.Exit(1)
		}
	}

	adminClient.Close()

}

// func kafkaConsumer() {
// 	//create consumer
// 	topicRespon := "user-respon"
// 	var topics []string
// 	topics = append(topics, conf.GetTopicRequest())
// 	sigchan := make(chan os.Signal, 1)
// 	signal.Notify(sigchan, syscall.SIGINT, syscall.SIGTERM)
// 	//earliest: automatically reset the offset to the earliest offset
// 	//latest: automatically reset the offset to the latest offset
// 	//none: throw exception to the consumer if no previous offset is found for the consumer's group
// 	c, err := kafka.NewConsumer(&kafka.ConfigMap{
// 		"bootstrap.servers":               conf.GetBroker(),
// 		"group.id":                        conf.GetGroup(),
// 		"session.timeout.ms":              6000,
// 		"go.events.channel.enable":        true,
// 		"go.application.rebalance.enable": true,
// 		"auto.offset.reset":               "latest",
// 	})
// 	if err != nil {
// 		log.Printf("Failed to create consumer: %s\n", err)
// 		os.Exit(1)
// 	}
// 	log.Printf("Created Consumer %v\n", c)
// 	err = c.SubscribeTopics(topics, nil)
// 	run := true
// 	for run == true {
// 		select {
// 		case sig := <-sigchan:
// 			log.Printf("Caught signal %v: terminating\n", sig)
// 			run = false
// 		case ev := <-c.Events():
// 			switch e := ev.(type) {
// 			case kafka.AssignedPartitions:
// 				log.Printf("%% %v\n", e)
// 				c.Assign(e.Partitions)
// 			case kafka.RevokedPartitions:
// 				log.Printf("%% %v\n", e)
// 				c.Unassign()
// 			case *kafka.Message:
// 				//log.Printf("%% Message on %s:\n%s\n",e.TopicPartition, string(e.Value))
// 				s := string(e.Value)
// 				// data := make(map[string]interface{})
// 				data := &Request{}
// 				json.Unmarshal([]byte(s), &data)

// 				log.Printf("s : %v", s)
// 				log.Printf("data : %s", data)
// 				resp := message(http.StatusOK, true, "success")
// 				// keys := reflect.ValueOf(data).MapKeys()
// 				// log.Println(keys)
// 				resp["request"] = data.Request
// 				resp["payload"] = data.Payload
// 				resp["token"] = data.Token
// 				keyId := uuid.Must(uuid.NewV4(), err).String()
// 				var dataHasil = map[string]string{"company": "34", "name": "Dapur Cokelat", "address": "Jl. Buaran", "key": keyId}
// 				resp["output"] = dataHasil

// 				// jsonString, err := json.Marshal(resp)
// 				// log.Printf("jsonString : %v", jsonString)

// 				log.Printf("LANJUT PROSES UNTUK KE DATA")
// 				err = sendToBroker(topicRespon, resp, "test", "test")
// 				// err = sendToBroker(procdr, topicRespon, resp, "test", "test")
// 				if err != nil {
// 					// handle this errors
// 					log.Printf("Error sendToBroker %v\n", err)
// 				}

// 			case kafka.PartitionEOF:
// 				log.Printf("%% Reached %v\n", e)
// 			case kafka.Error:
// 				log.Printf("Error: %v\n", e)
// 				run = false

// 			}
// 			log.Printf("cek...\n")
// 		}
// 		log.Printf("cek-2...\n")
// 	}

// 	log.Printf("Closing consumer\n")
// 	log.Printf("End at %v\n", time.Now())
// 	c.Close()
// 	os.Exit(3)
// }
