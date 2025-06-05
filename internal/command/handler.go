package command

import (
	"fmt"
	"strings"
	"sync"
)

type CommandFunc func(args []string) string

type commandHandler struct {
	commands map[string]CommandFunc
	mu       sync.RWMutex
}

var defaultCommandHandler = &commandHandler{
	commands: make(map[string]CommandFunc),
}

// 注册命令
func Register(cmd string, fn CommandFunc) {
	defaultCommandHandler.mu.Lock()
	defaultCommandHandler.commands[cmd] = fn
	defaultCommandHandler.mu.Unlock()
}

// 处理命令
func Handle(input string) string {
	parts := strings.Fields(input)
	if len(parts) == 0 {
		return "无命令输入"
	}
	cmd := parts[0]
	args := parts[1:]
	defaultCommandHandler.mu.RLock()
	fn := defaultCommandHandler.commands[cmd]
	defaultCommandHandler.mu.RUnlock()
	if fn == nil {
		return fmt.Sprintf("未知命令: %s", cmd)
	}
	return fn(args)
}

func ListCommands() []string {
	defaultCommandHandler.mu.RLock()
	defer defaultCommandHandler.mu.RUnlock()
	cmds := make([]string, 0, len(defaultCommandHandler.commands))
	for cmd := range defaultCommandHandler.commands {
		cmds = append(cmds, cmd)
	}
	return cmds
}