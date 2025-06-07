package command

import (
	"fmt"
)

// helpHandler implements the Handler interface for the "help" command
type helpHandler struct{}

func (h *helpHandler) Execute(args []string) error {
	fmt.Println("Available commands:")
	for _, cmd := range ListCommands() {
		fmt.Println(" -", cmd)
	}
	return nil
}

func init() {
	RegisterHandler("help", &helpHandler{})
}