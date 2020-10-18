package main

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/ChicoCodes/twitchbot/commands"
	"github.com/ChicoCodes/twitchbot/messages"
)

// matchReply creates a subscriber that replies with the given text if the
// incoming message matches the provided regular expression.
func matchReply(re *regexp.Regexp, response string) messages.Notify {
	return func(notification messages.Notification) {
		if re.MatchString(strings.ToLower(notification.Message.Text)) {
			notification.Reply(response)
		}
	}
}

func registerDefaultSubscribers(producer *messages.Producer) error {
	commandsManager, err := commands.New()
	if err != nil {
		return fmt.Errorf("failed to create the commands manager: %w", err)
	}
	subscribers := []messages.Notify{
		commandsManager.Subscribe,
		// replies with :w whenever someone says 'salve'.
		matchReply(regexp.MustCompile(`salve`), "/me :w"),
		// replies with a Vim error message whenever someone says ':qa'.
		matchReply(regexp.MustCompile(`^:qa$`), `/me E162: No write since last change for buffer "chat"`),
		// logs all messages to stdout.
		func(notification messages.Notification) {
			msg := notification.Message
			fmt.Printf("[%s] %s: %s\n", msg.Timestamp, msg.User.DisplayName, msg.Text)
		},
	}

	for _, subscriber := range subscribers {
		_, err := producer.Subscribe(subscriber)
		if err != nil {
			return fmt.Errorf("failed to register subscriber: %w", err)
		}
	}

	return nil
}
