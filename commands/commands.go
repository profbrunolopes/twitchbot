package commands

import (
	"fmt"
	"io/ioutil"
	"math/rand"
	"strings"
	"time"

	"github.com/ChicoCodes/twitchbot/messages"
)

// Command represents a bot command.
type Command struct {
	MinArgs int
	Exec    func(args []string, notification *messages.Notification)
	Help    string
}

func commandFn(exec func([]string, *messages.Notification)) *Command {
	return &Command{Exec: exec}
}

type Commands struct {
	commands map[string]*Command
}

// New creates the command processor.
func New() (*Commands, error) {
	rand.Seed(time.Now().Unix())
	commands := map[string]*Command{
		"say": startSaying(),
		"sorteio": commandFn(func(_ []string, notification *messages.Notification) {
			notification.Reply(fmt.Sprintf("parabéns %s, você ganhou uma licença do Vim!", notification.Message.User))
		}),
		"ban": {
			MinArgs: 1,
			Exec: func(args []string, notification *messages.Notification) {
				reasons := []string{
					"chamar biscoito de bolacha",
					"usar dark mode",
					"usar emacs",
					"comprar rodinha da Apple de 4 mil reais",
					"comprar um Mac Pro pra jogar Minecraft",
					"usar VSCode",
					"jogar Tibia",
					"usar Windows",
					"abusar do !say",
					"depositar 89 mil reais na conta da Micheque Kappa",
					"defender C",
					"não gostar do @pokemaobr",
				}
				reason := reasons[rand.Intn(len(reasons))]
				notification.Reply(fmt.Sprintf("/me %s baniu %s por %s", notification.Message.User, strings.Join(args, " "), reason))
			},
		},
		"colorscheme": {
			MinArgs: 1,
			Exec: func(args []string, notification *messages.Notification) {
				const file = "/tmp/colorscheme.txt"
				err := ioutil.WriteFile(file, []byte(strings.Join(args, " ")), 0600)
				if err != nil {
					notification.Reply(fmt.Sprintf("deu erro escrevendo o arquivo com colorscheme %v", err))
					return
				}
			},
		},
	}
	return &Commands{commands: commands}, nil
}

func (c *Commands) Subscribe(notification messages.Notification) {
	text := notification.Message.Text
	if !strings.HasPrefix(text, "!") {
		return
	}

	parts := strings.Split(text, " ")
	commandName := parts[0][1:]
	args := parts[1:]
	if command := c.commands[commandName]; command != nil {
		if len(args) < command.MinArgs {
			notification.Reply(command.Help)
			return
		}
		command.Exec(args, &notification)
	}
}
