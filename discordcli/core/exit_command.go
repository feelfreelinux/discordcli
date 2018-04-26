package core

import (
	"os"
)

type exitCommand struct{}

func (exitCommand) Execute(args []string) error {
	os.Exit(0)

	return nil
}

func (exitCommand) Name() string {
	return "exit"
}

func (exitCommand) HelpText() string {
	return "/exit - exits the application"
}

func init() {
	RegisteredCommands = append(RegisteredCommands, exitCommand{})
}
