package commands

import (
	"log"
	"os/exec"
	"strings"

	"github.com/ChicoCodes/twitchbot/messages"
)

func startSaying() *Command {
	const bufferLen = 50
	ch := make(chan string, bufferLen)
	go func() {
		for message := range ch {
			err := exec.Command("say", "-v", "Joana", message).Run()
			if err != nil {
				log.Printf("failed to run command: %v", err)
			}
		}
	}()
	return &Command{
		MinArgs: 1,
		Help:    "diz sua mensagem aí pow",
		Exec: func(args []string, _ *messages.Notification) {
			ch <- strings.Join(args, " ")
		},
	}
}
