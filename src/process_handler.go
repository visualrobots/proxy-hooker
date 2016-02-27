package main

import (
	"log"
	"os/exec"
	"strings"
)

type ProcessHandler struct {
	Command string
}

func (p *ProcessHandler) Reload() {
	go func() {
		args := strings.Split(p.Command, " ")
		cmd := exec.Command(args[0], args[1:]...)
		err := cmd.Run()
		if err != nil {
			log.Fatal(err)
		}

		log.Printf("Command '%s' executed", p.Command)
	}()
}

func NewProcessHandler(command string) *ProcessHandler {
	return &ProcessHandler{Command: command}
}
