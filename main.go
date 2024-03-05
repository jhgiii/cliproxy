package main

import (
	"fmt"
	prompt "github.com/c-bata/go-prompt"
	"github.com/jhgiii/cliproxy/sshclient"
	"golang.org/x/crypto/ssh"
	"os"
	"strings"
)

func executor(t string) {
	cmd := strings.Split(t, " ")
	if cmd[0] == "connect" {
		err := sshclient.Connect(ssh.ClientConfig{}, cmd[1])
		if err != nil {
			fmt.Println("Error connecting to " + cmd[1])
		}
	}
	if t == "exit" {
		os.Exit(0)
	}
}

func completer(t prompt.Document) []prompt.Suggest {
	return []prompt.Suggest{
		{Text: "connect"},
		{Text: "exit"},
	}
}

func main() {
	p := prompt.New(
		executor,
		completer,
	)
	p.Run()
}
