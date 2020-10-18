package main

import (
	"fmt"
	"regexp"
	"strings"

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

var defaultSubscribers = []messages.Notify{
	salve,
	qa,
	logger,
}

// replies with :w whenever someone says 'salve'
var salve = matchReply(regexp.MustCompile(`salve`), "/me :w")

// replies with a Vim error message whenever someone says ':qa'
var qa = matchReply(regexp.MustCompile(`^:qa$`), `/me E162: No write since last change for buffer "chat"`)

// logs all messages to stdout
var logger = func(notification messages.Notification) {
	msg := notification.Message
	fmt.Printf("[%s] %s: %s\n", msg.Timestamp, msg.User.DisplayName, msg.Text)
}
