package commands

import (
	"fmt"
	"io/ioutil"
	"math/rand"
	"sort"
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
	aliases  map[string]string
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
					"abusar do !say",
					"depositar 89 mil reais na conta da Micheque Kappa",
					"defender C",
					"não gostar do @pokemaobr",
					"colocar purê no cachorro quente",
				}
				//nolint:gosec
				reason := reasons[rand.Intn(len(reasons))]
				target := strings.Join(args, " ")
				msg := fmt.Sprintf("/me %s baniu %s por %s", notification.Message.User, target, reason)
				notification.Reply(msg)
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
		"music": commandFn(music),
	}
	aliases := map[string]string{
		"musica":     "music",
		"tocando":    "music",
		"nowplaying": "music",
		"comandos":   "commands",
	}
	c := Commands{commands: commands, aliases: aliases}
	c.commands["commands"] = commandFn(c.displayCommands)
	return &c, nil
}

func (c *Commands) Subscribe(notification messages.Notification) {
	text := notification.Message.Text
	if !strings.HasPrefix(text, "!") {
		return
	}

	parts := strings.Split(text, " ")
	commandName := parts[0][1:]
	if aliased, ok := c.aliases[commandName]; ok {
		commandName = aliased
	}
	args := parts[1:]
	if command := c.commands[commandName]; command != nil {
		if len(args) < command.MinArgs {
			if command.Help != "" {
				notification.Reply(command.Help)
			}
			return
		}
		command.Exec(args, &notification)
	}
}

func (c *Commands) displayCommands(_ []string, notification *messages.Notification) {
	commands := make([]string, 0, len(c.commands))
	for command := range c.commands {
		commands = append(commands, "!"+command)
	}
	sort.Strings(commands)
	notification.Reply(fmt.Sprintf("/me available commands: %s", strings.Join(commands, ", ")))
}
