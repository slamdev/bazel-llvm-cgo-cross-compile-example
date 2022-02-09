package main

import (
	"github.com/confluentinc/confluent-kafka-go/kafka"
)

func main() {
	print(kafka.ErrMaxPollExceeded)
}
