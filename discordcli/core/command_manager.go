package core

import (
	"strings"
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
func (CommandManager) ExecuteCommand(msg string) {
	cmd := strings.TrimLeft(msg, "/")
	args := strings.Split(cmd, " ")
	for i := range args {
		if RegisteredCommands[i].Name() == args[0] {
			RegisteredCommands[i].Execute(args[1:])
		}
	}
	panic("COMMAND NOT FOUND")
}
