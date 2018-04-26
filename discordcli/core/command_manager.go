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
ExecuteCommand finds a command to run and runs
*/
func (CommandManager) ExecuteCommand(msg string) error {
	cmd := strings.TrimLeft(msg, "/")
	args := strings.Split(cmd, " ")

	for i := range args {
		if RegisteredCommands[i].Name() == strings.TrimSpace(args[0]) {

			return RegisteredCommands[i].Execute(args[1:])

		}
	}

	return errors.New("Command " + args[0] + "not found")

}
