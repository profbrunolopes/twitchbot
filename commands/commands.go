package commands

import (
	"fmt"
	"io/ioutil"
	"strings"
	"time"

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

// sistema de permissão melhor:
//
// -> por comando (checa se o usuário é ChicoCodes ou chicocodesbot - ou
//    nightbot)
//      sugestão do CodeShow: usar palavras proibidas
//      solução ideal: integrar com a API do Twitch (QUAL API?!)
//
// -> integrado com o twitch pra saber role do usuário
//    + comandos só pra mod (& broadcaster)
//    + comandos pra todo mundo
//
// -> colorscheme pode ir pra lojinha, e o bot executa o comando.

func (c *Commands) Subscribe(notification messages.Notification) {
	if strings.HasPrefix(notification.Message.Text, "!starttimer") {
		// sistema de permissão avançado
		if strings.ToLower(notification.Message.User) != "chicocodes" {
			return
		}

		parts := strings.SplitN(notification.Message.Text, " ", 3)
		if len(parts) != 3 {
			notification.Reply("comando inválido")
			return
		}
		duration, err := time.ParseDuration(parts[1])
		if err != nil {
			notification.Reply(fmt.Sprintf("%v não é uma duração válida", parts[1]))
			return
		}
		go func() {
			<-time.After(duration)
			notification.Reply(fmt.Sprintf("/me timer %q encerrado", parts[2]))
		}()
	}

	if strings.HasPrefix(notification.Message.Text, "!ban") {
		parts := strings.SplitN(notification.Message.Text, " ", 2)
		if len(parts) != 2 {
			notification.Reply("comando inválido")
			return
		}
		notification.Reply(fmt.Sprintf("/me baniu %s por usar dark mode", parts[1]))
	}

	if strings.HasPrefix(notification.Message.Text, "!colorscheme") {
		parts := strings.SplitN(notification.Message.Text, " ", 2)
		if len(parts) != 2 {
			notification.Reply("comando inválido")
			return
		}
		const file = "/tmp/colorscheme.txt"
		err := ioutil.WriteFile(file, []byte(parts[1]), 0644)
		if err != nil {
			notification.Reply(fmt.Sprintf("deu erro escrevendo o arquivo com colorscheme %v", err))
			return
		}
	}
}
