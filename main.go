package main

import (
	"flag"
	"fmt"
	"log"
	"strings"

	"github.com/ChicoCodes/twitchbot/commands"
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
	_, err = producer.Subscribe(func(notification messages.Notification) {
		msg := notification.Message
		fmt.Printf("[%s] %s: %s\n", msg.Timestamp, msg.User, msg.Text)
	})
	if err != nil {
		log.Fatal(err)
	}
	_, err = producer.Subscribe(func(notification messages.Notification) {
		if strings.Contains(strings.ToLower(notification.Message.Text), "salve") {
			notification.Reply("/me :w")
		}
	})
	if err != nil {
		log.Fatal(err)
	}
	_, err = producer.Subscribe(func(notification messages.Notification) {
		if strings.ToLower(notification.Message.Text) == ":qa" {
			notification.Reply(`/me E162: No write since last change for buffer "chat"`)
		}
	})
	if err != nil {
		log.Fatal(err)
	}
	cmds, err := commands.New()
	if err != nil {
		log.Fatal(err)
	}
	_, err = producer.Subscribe(cmds.Subscribe)
	if err != nil {
		log.Fatal(err)
	}
	err = producer.Start()
	if err != nil {
		log.Fatal(err)
	}
}
