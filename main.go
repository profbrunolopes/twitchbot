package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/ChicoCodes/twitchbot/messages"
	"github.com/kelseyhightower/envconfig"
)

func main() {
	channel := flag.String("channel", "", "channel to connect to")
	flag.Parse()

	options := messages.ProducerOptions{Channel: *channel}
	err := envconfig.Process("twitch", &options)
	if err != nil {
		log.Fatal(err)
	}

	producer := messages.NewProducer(&options)
	producer.Subscribe(func(msg messages.Message) {
		fmt.Printf("[%s] %s: %s\n", msg.Timestamp, msg.User, msg.Text)
	})
	err = producer.Start()
	if err != nil {
		log.Fatal(err)
	}
}
