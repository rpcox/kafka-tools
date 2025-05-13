package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
)

const _tool = "kppipe"
const _version = "0.1.0"

// Listen to all events on the default events channel
func EventListener(p *kafka.Producer, done chan interface{}) {

	for {
		select {
		case <-done:
			//fmt.Println("exit event listener")
			return
		default:
			for e := range p.Events() {
				switch ev := e.(type) {
				case *kafka.Message:
					if ev.TopicPartition.Error != nil {
						fmt.Fprintf(os.Stderr, "publish fail: %v\n", ev.TopicPartition.Error)
					} else {
						fmt.Fprintf(os.Stdout, "publish success: topic:%s in partition[%d] @ offset %v\n",
							*ev.TopicPartition.Topic, ev.TopicPartition.Partition, ev.TopicPartition.Offset)
					}
				case kafka.Error:
					fmt.Fprintf(os.Stderr, "%%error: %v\n", ev)
				default:
					fmt.Fprintf(os.Stderr, "%%unhandled: %s\n", ev)
				}
			}
		}
	}

}

func PipedInput() bool {
	fi, err := os.Stdin.Stat()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(ExitUnableToStatStdin)
	}

	if (fi.Mode() & os.ModeCharDevice) == 0 {
		return true
	}

	return false
}

func main() {
	_broker := flag.String("b", "", "Override the config with a specific broker / bootstrap server")
	_config := flag.String("p", "", "Location of Kafka producer properties file")
	_keygen := flag.String("keygen", "uuid", "Location of Kafka producer properties file")
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

	done := make(chan interface{})
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
	go SigHandler(sig, done)

	go EventListener(p, done)

	if !PipedInput() {
		fmt.Fprintf(os.Stderr, "Input not detected on stdin pipe\n")
		os.Exit(ExitInputNotOnPipe)
	}

	kg := GetKeyGenerator(*_keygen)
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		if scanner.Err() == io.EOF {
			fmt.Println("end of pipe")
			break
		}

		err = p.Produce(&kafka.Message{
			TopicPartition: kafka.TopicPartition{Topic: _topic, Partition: kafka.PartitionAny},
			Key:            []byte(kg.KeyGen()),
			Value:          scanner.Bytes(),
		}, nil)

		if err != nil {
			if err.(kafka.Error).Code() == kafka.ErrQueueFull {
				time.Sleep(time.Second)
				continue
			}
			fmt.Printf("fail to produce: %v\n", err)
		}
	}

	for n := p.Flush(1000); n > 0; {
		if n%10 == 0 {
			fmt.Fprintf(os.Stderr, "flushing...\n")
		}
	}

	close(done)
}
