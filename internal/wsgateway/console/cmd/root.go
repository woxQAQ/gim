package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/chzyer/readline"
	"github.com/spf13/cobra"

	"github.com/woxQAQ/gim/internal/wsgateway"
	"github.com/woxQAQ/gim/pkg/logger"
)

var (
	gateway wsgateway.Gateway
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Short: "WebSocket Gateway Console",
	// 禁用自动添加的flags
	DisableFlagParsing: true,
	// 禁用completion命令
	CompletionOptions: cobra.CompletionOptions{
		DisableDefaultCmd: true,
	},
}

// getCompletions 根据当前输入返回可能的补全选项
func getCompletions(line string) []string {
	// 获取当前输入的所有部分
	parts := strings.Fields(line)

	// 如果没有输入，返回所有可用的命令
	if len(parts) == 0 {
		commands := []string{"broadcast", "check", "stats", "help", "exit", "quit"}
		return commands
	}

	// 如果正在输入第一个词，过滤可用的命令
	if len(parts) == 1 {
		commands := []string{"broadcast", "check", "stats", "help", "exit", "quit"}
		matches := []string{}
		for _, cmd := range commands {
			if strings.HasPrefix(cmd, parts[0]) {
				matches = append(matches, cmd)
			}
		}
		return matches
	}

	// 根据不同的命令提供参数补全
	switch parts[0] {
	case "check":
		// check命令需要一个user_id参数
		if len(parts) == 2 {
			// 这里可以添加用户ID的补全逻辑
			return []string{}
		}
	case "broadcast":
		// broadcast命令接受任意文本消息
		return []string{}
	}

	return []string{}
}

// Execute 执行根命令
func Execute(g wsgateway.Gateway, l logger.Logger) {
	gateway = g

	// 初始化readline实例
	rl, err := readline.NewEx(&readline.Config{
		Prompt:            "console> ",
		HistoryFile:       "/tmp/gim_console.history",
		HistoryLimit:      100,
		HistorySearchFold: true,
		AutoComplete:      readline.NewPrefixCompleter(),
	})
	if err != nil {
		fmt.Printf("初始化控制台失败: %v\n", err)
		os.Exit(1)
	}
	defer rl.Close()

	// 设置自定义的补全函数
	rl.Config.AutoComplete = readline.NewPrefixCompleter(
		readline.PcItemDynamic(func(line string) []string {
			return getCompletions(line)
		}),
	)

	// 设置根命令的Run函数，使其进入交互式模式
	rootCmd.Run = func(cmd *cobra.Command, args []string) {
		fmt.Println("WebSocket Gateway Console 已启动，输入 help 查看可用命令")
		for {
			// 读取用户输入
			line, err := rl.Readline()
			if err != nil {
				if err == readline.ErrInterrupt {
					continue
				}
				break
			}

			// 处理空输入
			if line = strings.TrimSpace(line); line == "" {
				continue
			}

			// 处理退出命令
			if line == "exit" || line == "quit" {
				fmt.Println("正在关闭网关服务器...")
				// 关闭网关服务
				if err := gateway.Stop(); err != nil {
					fmt.Printf("关闭网关服务失败: %v\n", err)
				}
				os.Exit(0)
			}

			// 执行命令
			rootCmd.SetArgs(strings.Fields(line))
			l.Disable()
			err = rootCmd.Execute()
			l.Enable()
			if err != nil {
				fmt.Printf("执行命令失败: %v\n", err)
			}
			l.Enable()
		}
	}

	// 执行根命令
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
