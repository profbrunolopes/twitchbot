package main

import (
	"flag"
	"fmt"

	"github.com/ChicoCodes/twitchbot/messages"
)

func main() {
	channel := flag.String("channel", "ChicoCodes", "channel to listen to")
	flag.Parse()

	producer := messages.NewProducer(*channel)
	messages, err := producer.Subscribe("random string")
	if err != nil {
		panic(err)
	}
	for message := range messages {
		fmt.Printf("[%s] %s: %s\n", message.Timestamp, message.User, message.Text)
	}
}
