package commands

import (
	"github.com/ChicoCodes/twitchbot/messages"
)

type Commands struct {
}

func New() (*Commands, error) {
	return &Commands{}, nil
}

// !addcommand
// !rmcommand
// !commands

// syntax for addcommand:
//
//   - !addcommand <comando> <ação> [parâmetros para a ação ...]
//
//   exemplos de ações: alias, write-file, reply,
//   exemplos de ações: alias, write-file, reply,
//   exemplos de ações: alias, write-file, reply,
//   exemplos de ações: alias, write-file, reply,
//   exemplos de ações: alias, write-file, reply,
//
// !addcommand comandos alias commands
//
// !comandos => !commands
//
// !addcommand color write-file /tmp/colorscheme.txt

// !color gruvbox dark
//
// -> escreve /tmp/colorscheme.txt
// -> vamos criar um plugin genérico pro vim que monitora arquivos e chama uma
//    função sempre que o arquivo muda, passando o path como parâmetro da função
// -> vamos usar esse plugin para escrever uma função que lê o conteúdo do
// arquivo e define o colorscheme e background

// !ban

func (c *Commands) Subscribe(notification messages.Notification) {
}
