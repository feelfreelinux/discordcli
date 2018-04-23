package core

type exitCommand struct{}

func (exitCommand) Execute(args []string) error {
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
