package main

import (
	"bytes"
	"fmt"
	"log"

	"golang.org/x/crypto/ssh"
)

func main() {
	const (
		ip   = "192.168.100.1:22"
		user = "user"
	)
	config := &ssh.ClientConfig{
		User: user,
		Auth: []ssh.AuthMethod{
			ssh.Password("secrets"),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(), // password認証は設定
	}
	client, err := ssh.Dial("tcp", ip, config)
	if err != nil {
		log.Fatal("Failed to dial: ", err)
	}
	defer client.Close()

	session, err := client.NewSession()
	if err != nil {
		log.Fatal("Failed to create session: ", err)
	}
	defer session.Close()

	modes := ssh.TerminalModes{
		ssh.ECHO: 0, // supress echo
	}

	// run terminal session
	if err := session.RequestPty("xterm", 50, 80, modes); err != nil {
		log.Fatal(err)
	}

	// start remote shell
	if err := session.Shell(); err != nil {
		log.Fatal(err)
	}

	// TODO: open interactive session

	var b bytes.Buffer
	session.Stdout = &b
	if err := session.Run("ls"); err != nil {
		log.Fatal("Failed to run: " + err.Error())
	}
	fmt.Println(b.String())
	fmt.Println("Hello World")
}
