package main

import (
	"bufio"
	"crypto/rand"
	"fmt"
	"log"
	"os"
	"runtime/debug"
	"strings"
	"syscall"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
)

var (
	_commit string
	_branch string

	ExitNewProducerFail   = 2
	ExitUnableToStatStdin = 3
	ExitInputNotOnPipe    = 4
)

func Version(b bool) {
	if b {
		fmt.Fprintf(os.Stderr, "%s v%s\n", _tool, _version)
		if _commit != "" && _branch != "" {
			// go build -ldflags="-X main._commit=$(git rev-parse --short HEAD) -X main._branch=$(git branch | awk '{print $2}')"
			fmt.Fprintf(os.Stderr, "commit: %s, branch: %s\n", _commit, _branch)
		}

		if info, ok := debug.ReadBuildInfo(); ok {
			fmt.Fprintf(os.Stderr, "go_version %s\n", info.GoVersion)
		}

		os.Exit(0)
	}
}

// Read a key=value file for Kafka properties
// Comments in the file are marked with a '#' in column 0
func ReadConfig(configFile string) kafka.ConfigMap {

	m := make(map[string]kafka.ConfigValue)

	file, err := os.Open(configFile)
	if err != nil {
		log.Fatalf("config open: %s", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if !strings.HasPrefix(line, "#") && len(line) != 0 {
			kv := strings.Split(line, "=")
			property := strings.TrimSpace(kv[0])
			value := strings.TrimSpace(kv[1])
			m[property] = value
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatal("config read: %s", err)
	}

	return m
}

func SigHandler(sig chan os.Signal, done chan interface{}) {
	for {
		signal := <-sig
		if signal == syscall.SIGINT || signal == syscall.SIGTERM {
			fmt.Println("exit signal handler")
			done <- true
			break
		}

	}

	log.Println("signal: exiting")
}

func RandString(length int) (string, error) {
	const chars = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	bytes := make([]byte, length)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}

	for i, b := range bytes {
		bytes[i] = chars[b%byte(len(chars))]
	}

	return string(bytes), nil
}

// SDG
