package command

import (
	"fmt"
	"strings"
)

// Handler interface that all command handlers must implement
type Handler interface {
	Execute(args []string) error
}

// Registry to hold all handlers
var handlers = make(map[string]Handler)

// RegisterHandler registers a handler for a command
func RegisterHandler(command string, handler Handler) {
	handlers[strings.ToLower(command)] = handler
}

// ParseAndExecute parses the input and executes the corresponding handler
// Command format: "help arg1 arg2 ..."
func ParseAndExecute(input string) error {
	parts := strings.Fields(input)
	if len(parts) == 0 {
		return fmt.Errorf("no command provided")
	}
	cmd := strings.ToLower(parts[0])
	handler, ok := handlers[cmd]
	if !ok {
		return fmt.Errorf("unknown command: %s", cmd)
	}
	// Pass all arguments after the command as args (can be multiple)
	return handler.Execute(parts[1:])
}

func ListCommands() []string {
	cmds := make([]string, 0, len(handlers))
	for cmd := range handlers {
		cmds = append(cmds, cmd)
	}
	return cmds
}