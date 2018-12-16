package common

import "github.com/urfave/cli"

var commands []cli.Command

// RegisterCommand append the command to command array.
func RegisterCommand(command cli.Command) {
	commands = append(commands, command)
}

// GetCommands return the registered commands.
func GetCommands() []cli.Command {
	return commands
}
