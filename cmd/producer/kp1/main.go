package main

import (
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
			return
		default:
			for e := range p.Events() {
				switch ev := e.(type) {
				case *kafka.Message:
					msg <- ev
				case kafka.Error:
					fmt.Fprintf(os.Stderr, "%%error: %v\n", ev)
				default:
					fmt.Fprintf(os.Stderr, "%%unhandled: %s\n", ev)
				}
			}
		}
	}

}

func Flush(p *kafka.Producer) {
	for n := p.Flush(1000); n > 0; {
		if n%100 == 0 {
			fmt.Fprintf(os.Stderr, "flushing...\n")
		}
	}
}

func main() {
	_broker := flag.String("b", "", "Identify a specific broker / bootstrap server")
	_config := flag.String("p", "", "Location of Kafka producer properties file")
	_count := flag.Int("count", 5, "Specify the number of messages to publish")
	_length := flag.Int("l", 64, "Specify the length of the random message")
	_topic := flag.String("t", "", "Identify the Kafka topic")
	_version := flag.Bool("version", false, "Display version and exit")
	flag.Parse()

	Version(*_version)

	config := ReadConfig(*_config)
	if *_broker != "" {
		config["bootstrap.servers"] = *_broker
	}

	p, err := kafka.NewProducer(&config)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(ExitNewProducerFail)
	}
	defer p.Close()
	defer Flush(p)

	msg := make(chan *kafka.Message)
	defer close(msg)

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
			return
		}
	}

}
