package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
)

const _tool = "kp"
const _version = "0.1.0"

// Listen to all events on the default events channel
func EventListener(p *kafka.Producer, done chan interface{}, msg chan *kafka.Message, wg *sync.WaitGroup) {
	defer wg.Done()

	for {
		select {
		case <-done:
			fmt.Println("EventListener: exiting")
			return
		case e := <-p.Events():
			//default:
			//for e := range p.Events() {
			switch ev := e.(type) {
			case *kafka.Message:
				msg <- ev
			case kafka.Error:
				fmt.Fprintf(os.Stderr, "%%error: %v\n", ev)
			default:
				if ev != nil {
					fmt.Fprintf(os.Stderr, "%%unhandled: %v\n", ev)
				}
			}
			//}
		}
	}

}

func Flush(p *kafka.Producer) {
	for p.Flush(1000) > 0 {
		fmt.Fprintf(os.Stderr, "flushing...\n")
	}
}

func UseSpecificBroker(b bool, kcm kafka.ConfigMap, broker string) {
	if b {
		kcm["bootstrap.servers"] = broker
	}
}

func TopicExist(kcm *kafka.ConfigMap, topic *string) bool {
	ac, err := kafka.NewAdminClient(kcm)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(ExitNewAdminClientFail)
	}
	defer ac.Close()

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()
	TopicNames := []string{*topic}
	describeTopicsResult, err := ac.DescribeTopics(
		ctx, kafka.NewTopicCollectionOfTopicNames(TopicNames),
		//	kafka.SetAdminOptionIncludeAuthorizedOperations(
		//		includeAuthorizedOperations)
	)
	if err != nil {
		fmt.Printf("failed to describe topics: %s\n", err)
		os.Exit(1)
	}

	for _, t := range describeTopicsResult.TopicDescriptions {
		if t.Name == *topic {
			if t.Error.Code() == kafka.ErrUnknownTopicOrPart {
				fmt.Printf("topic: %s has error(%d): %s\n", t.Name, t.Error.Code(), t.Error)
				return false
			}

			return true
		}
	}

	fmt.Printf("topic '%s' does NOT exist (not found in metadata).\n", topic)
	return false
}

func main() {
	_broker := flag.String("b", "", "Identify a specific broker / bootstrap server. e.g., 'spud:9092'")
	_config := flag.String("p", "", "Location of Kafka producer properties file")
	_count := flag.Int("count", 5, "Specify the number of messages to publish")
	_length := flag.Int("l", 64, "Specify the length of the random message")
	_topic := flag.String("t", "", "Identify the Kafka topic")
	_version := flag.Bool("version", false, "Display version and exit")
	flag.Parse()

	Version(*_version)
	config := ReadConfig(*_config)
	UseSpecificBroker(*_broker != "", config, *_broker)

	if !TopicExist(&config, _topic) {
		os.Exit(ExitTopicNotExist)
	}

	p, err := kafka.NewProducer(&config)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(ExitNewProducerFail)
	}
	defer p.Close()
	msg := make(chan *kafka.Message)
	defer close(msg)

	//defer p.Flush(100)
	//	defer Flush(p)

	done := make(chan interface{})
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
	go SigHandler(sig, done)

	var wg sync.WaitGroup
	wg.Add(1)
	go EventListener(p, done, msg, &wg)

	msgCount := 0
	for msgCount < *_count {
		rs, _ := RandString(*_length)
		value := fmt.Sprintf("%d:%s", msgCount, rs)
		err = p.Produce(&kafka.Message{
			TopicPartition: kafka.TopicPartition{Topic: _topic, Partition: kafka.PartitionAny},
			Value:          []byte(value),
		}, nil)

		if err != nil {
			if err.(kafka.Error).Code() == kafka.ErrQueueFull {
				time.Sleep(time.Second)
				continue
			}
			fmt.Printf("%d: fail to produce: %v\n", msgCount, err)
		}
		msgCount++
	}

	for {
		select {
		case m := <-msg:
			if m.TopicPartition.Error != nil {
				fmt.Fprintf(os.Stderr, "publish fail: %v\n", m.TopicPartition.Error)
			} else {
				fmt.Fprintf(os.Stdout, "publish success: topic:%s in partition[%d] @ offset %v\n",
					*m.TopicPartition.Topic, m.TopicPartition.Partition, m.TopicPartition.Offset)
			}
			m = nil

			msgCount--
			if msgCount == 0 {
				return
			}
		case <-done:
			fmt.Println("main: exiting")
			return
		}
	}

}

// SDG
