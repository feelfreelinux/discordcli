package core

import (
	"strings"

	"github.com/pkg/errors"
)

/*
CommandManager processes and dispatches client commands
*/
type CommandManager struct {
}

/*
IsCommand checks whether a given message can be interpreted and exacuted as a client command
*/
func (CommandManager) IsCommand(msg string) bool {
	return strings.HasPrefix(msg, "/")
}

/*
ExecuteCommand finds a command to run and runs it
*/
func (CommandManager) ExecuteCommand(msg string) error {
	cmd := strings.TrimLeft(msg, "/")
	args := strings.Split(cmd, " ")
	commandName := strings.TrimSpace(args[0])
	for i := range args {
		if RegisteredCommands[i].Name() == commandName {

			return RegisteredCommands[i].Execute(args[1:])

		}
	}

	return errors.New("Command " + commandName + "not found")

}
