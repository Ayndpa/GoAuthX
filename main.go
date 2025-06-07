package main

import (
	"bufio"
	"fmt"
	"goauthx/internal/command"
	"goauthx/internal/web"
	"log"
	"os"
	"strings"
)

func main() {
	// 启动命令行监听协程
	go func() {
		reader := bufio.NewReader(os.Stdin)
		for {
			fmt.Print("> ")
			input, err := reader.ReadString('\n')
			if err != nil {
				log.Printf("读取命令失败: %v", err)
				continue
			}
			input = strings.TrimSpace(input)
			if input == "" {
				continue
			}
			if err := command.ParseAndExecute(input); err != nil {
				log.Printf("命令执行错误: %v", err)
			}
		}
	}()

	err := web.StartServer()
	if err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
