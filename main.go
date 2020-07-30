package main

import (
	"flag"
	"fmt"
	"log"
	"strings"

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
	producer.Subscribe(func(notification messages.Notification) {
		msg := notification.Message
		fmt.Printf("[%s] %s: %s\n", msg.Timestamp, msg.User, msg.Text)
	})
	producer.Subscribe(func(notification messages.Notification) {
		if strings.Contains(strings.ToLower(notification.Message.Text), "salve") {
			notification.Reply(":w")
		}
	})
	err = producer.Start()
	if err != nil {
		log.Fatal(err)
	}
}
