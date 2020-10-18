package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/ChicoCodes/twitchbot/commands"
	"github.com/ChicoCodes/twitchbot/messages"
	"github.com/kelseyhightower/envconfig"
)

func main() {
	channel := flag.String("channel", "", "channel to connect to")
	flag.Parse()

	if *channel == "" {
		flag.Usage()
		os.Exit(1)
	}

	options := messages.ProducerOptions{Channel: *channel}
	err := envconfig.Process("twitch", &options)
	if err != nil {
		log.Fatal(err)
	}

	producer := messages.NewProducer(&options)
	err = registerSubscribers(producer)
	if err != nil {
		log.Fatal(err)
	}
	err = producer.Start()
	if err != nil {
		log.Fatal(err)
	}
}

func registerSubscribers(producer *messages.Producer) error {
	for _, sub := range defaultSubscribers {
		_, err := producer.Subscribe(sub)
		if err != nil {
			return fmt.Errorf("failed to register default subscriber: %w", err)
		}
	}
	commandsSubscriber, err := commands.New()
	if err != nil {
		return fmt.Errorf("failed to create the commands manager: %w", err)
	}
	_, err = producer.Subscribe(commandsSubscriber.Subscribe)
	return err
}
