package core

/*
RegisteredCommands holds all commands registred by init() in command files
*/
var RegisteredCommands []ClientCommand

/*
ClientCommand provides information to the command manager about the arguments of the command.
*/
type ClientCommand interface {
	Name() string
	HelpText() string
	Execute([]string) error
}
