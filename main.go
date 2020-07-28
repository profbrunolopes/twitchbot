package main

import (
	"fmt"

	"github.com/ChicoCodes/twitchbot/messages"
)

func main() {
	producer := messages.NewProducer("ChicoCodes")

	messages, err := producer.Subscribe("random string")
	if err != nil {
		panic(err)
	}
	for message := range messages {
		fmt.Printf("[%s] %s: %s\n", message.Timestamp, message.User, message.Text)
	}
}
