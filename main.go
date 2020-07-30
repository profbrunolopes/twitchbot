package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/ChicoCodes/twitchbot/messages"
)

func main() {
	channel := flag.String("channel", "ChicoCodes", "channel to connect to")
	flag.Parse()

	producer := messages.NewProducer(*channel)
	producer.Subscribe(func(msg messages.Message) {
		fmt.Printf("[%s] %s: %s\n", msg.Timestamp, msg.User, msg.Text)
	})
	err := producer.Start()
	if err != nil {
		log.Fatal(err)
	}
}
