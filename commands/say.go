package commands

import (
	"log"
	"os/exec"
)

func startSaying() chan<- string {
	ch := make(chan string, 50)
	go func() {
		for message := range ch {
			err := exec.Command("say", "-v", "Joana", message).Run()
			if err != nil {
				log.Printf("failed to run command: %v", err)
			}
		}
	}()
	return ch
}
