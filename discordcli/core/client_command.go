package core

var RegisteredCommands []ClientCommand

/*
ClientCommand provides information to the command manager about the arguments of the command.
*/
type ClientCommand interface {
	Name() string
	HelpText() string
	Execute([]string) error
}
