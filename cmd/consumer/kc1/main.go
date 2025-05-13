// A simple consumer
package main

// consumer_example implements a consumer using the non-channel Poll() API
// to retrieve messages and events.

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
)

const _tool = "kc1"
const _version = "v0.1.0"

func main() {
	_broker := flag.String("b", "", "Identify a specific broker / bootstrap server")
	_config := flag.String("p", "", "Location of Kafka producer properties file")
	_topic := flag.String("t", "", "Identify the Kafka topic")
	_version := flag.Bool("version", false, "Display version and exit")
	flag.Parse()

	Version(*_version)

	config := ReadConfig(*_config)
	if *_broker != "" {
		config["bootstrap.servers"] = *_broker
	}
	c, err := kafka.NewConsumer(&config)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(ExitNewConsumerFail)
	}
	defer c.Close()

	done := make(chan interface{})
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	go SigHandler(sigChan, done)

	topics := []string{*_topic}
	err = c.SubscribeTopics(topics, nil)
	if err != nil {
		fmt.Fprintf(os.Stderr, "topic subscribe: %s\n", err)
		os.Exit(ExitTopicSubscribeFail)
	}

	for {
		select {
		case <-done:
			fmt.Fprintf(os.Stderr, "signal: exiting poll loop\n")
			return
		default:
			ev := c.Poll(100)
			if ev == nil {
				continue
			}

			switch e := ev.(type) {
			case *kafka.Message:
				// Do something with the record
				fmt.Fprintf(os.Stderr, "%% **\n  p=%v o=%v\n  %v\n  %s|%s\n",
					e.TopicPartition.Partition, e.TopicPartition.Offset, e.TopicPartition, string(e.Key), string(e.Value))
				if e.Headers != nil {
					fmt.Printf("%% header: %v\n", e.Headers)
				}

				_, err := c.StoreMessage(e)
				if err != nil {
					fmt.Fprintf(os.Stderr, "%% error: storage message: %s:\n", e.TopicPartition)
				}
			case kafka.Error:
				fmt.Fprintf(os.Stderr, "%% error: %v: %v\n", e.Code(), e)
				if e.Code() == kafka.ErrAllBrokersDown {
					fmt.Fprintf(os.Stderr, "%% exiting\n")
					return
				}

			case kafka.PartitionEOF:
				fmt.Println("partition EOF")
			case kafka.OffsetsCommitted:
				fmt.Printf("offset commit: %v\n", e)
			case nil:
				fmt.Println("% nil")
			default:
				/*if e == kafka.OffsetsCommitted {
					fmt.Printf("OffsetsCommitted: %v\n", e)
					continue
				}*/

				fmt.Printf("unhandled: %T => %v\n", e, e)
			}
		}
	}
}

// SDG
