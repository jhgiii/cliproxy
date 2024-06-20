package sshclient

import (
	"fmt"
	"log"
	"os"

	"golang.org/x/crypto/ssh"
	terminal "golang.org/x/term"
)

func Connect(config ssh.ClientConfig, server string) error {
	key, err := os.ReadFile("/Users/jim/.ssh/id_rsa")
	if err != nil {
		return err
	}

	signer, err := ssh.ParsePrivateKey(key)
	if err != nil {
		return err
	}
	config.User = "jim"
	config.Auth = []ssh.AuthMethod{
		ssh.PublicKeys(signer),
	}

	config.HostKeyCallback = ssh.InsecureIgnoreHostKey()

	fmt.Println(server)
	conn, err := ssh.Dial("tcp", server, &config)
	if err != nil {
		return err
	}
	defer conn.Close()

	session, err := conn.NewSession()
	if err != nil {
		return err
	}
	defer session.Close()

	// Set IO
	session.Stdout = os.Stdout
	session.Stderr = os.Stderr
	session.Stdin = os.Stdin

	// Set up terminal modes
	modes := ssh.TerminalModes{
		ssh.ECHO:          1, // enable echoing
		ssh.ECHOCTL:       1,
		ssh.TTY_OP_ISPEED: 14400, // input speed = 14.4kbaud
		ssh.TTY_OP_OSPEED: 14400, // output speed = 14.4kbaud
	}

	fileDescriptor := int(os.Stdin.Fd())
	if terminal.IsTerminal(fileDescriptor) {
		originalState, err := terminal.MakeRaw(fileDescriptor)
		if err != nil {
			log.Fatalf("Error at MakeRaw %v", err)
		}
		defer terminal.Restore(fileDescriptor, originalState)

		termWidth, termHeight, err := terminal.GetSize(fileDescriptor)
		if err != nil {
			log.Fatalf("Error at Setting Tterminal.GetSize: %v", err)
		}
		err = session.RequestPty("xterm-256color", termHeight, termWidth, modes)
		if err != nil {
			log.Fatalf("Error when requesting PTY: %v", err)
		}
	}
	err = session.Shell()
	if err != nil {
		log.Fatalf("Error when building shell: %v", err)
	}
	session.Wait()
	return nil
}
