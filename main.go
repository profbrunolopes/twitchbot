package main

import (
	"flag"
	"log"
	"os"

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
	err = registerDefaultSubscribers(producer)
	if err != nil {
		log.Fatal(err)
	}
	err = producer.Start()
	if err != nil {
		log.Fatal(err)
	}
}
