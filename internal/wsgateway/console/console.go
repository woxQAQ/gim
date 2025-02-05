package console

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/woxQAQ/gim/internal/types"
	"github.com/woxQAQ/gim/internal/wsgateway"
	"github.com/woxQAQ/gim/pkg/logger"
)

// Console 提供WebSocket网关的控制台界面
type Console struct {
	gateway wsgateway.Gateway
	logger  logger.Logger
	reader  *bufio.Reader
	running bool
}

// NewConsole 创建新的控制台实例
func NewConsole(gateway wsgateway.Gateway, logger logger.Logger) *Console {
	return &Console{
		gateway: gateway,
		logger:  logger,
		reader:  bufio.NewReader(os.Stdin),
	}
}

// Start 启动控制台
func (c *Console) Start() {
	c.running = true
	fmt.Println("WebSocket Gateway Console - Type 'help' for commands")

	for c.running {
		fmt.Print("> ")
		line, err := c.reader.ReadString('\n')
		if err != nil {
			c.logger.Error("Failed to read console input", logger.Error(err))
			continue
		}

		// 处理命令
		c.handleCommand(strings.TrimSpace(line))
	}
}

// Stop 停止控制台
func (c *Console) Stop() {
	c.running = false
}

// handleCommand 处理控制台命令
func (c *Console) handleCommand(cmd string) {
	args := strings.Fields(cmd)
	if len(args) == 0 {
		return
	}

	switch args[0] {
	case "help":
		c.printHelp()
	case "stats":
		c.printStats()
	case "broadcast":
		if len(args) < 2 {
			fmt.Println("Usage: broadcast <message>")
			return
		}
		c.broadcastMessage(strings.Join(args[1:], " "))
	case "check":
		if len(args) != 2 {
			fmt.Println("Usage: check <user_id>")
			return
		}
		c.checkUser(args[1])
	case "exit", "quit":
		c.Stop()
	default:
		fmt.Println("Unknown command. Type 'help' for available commands")
	}
}

// printHelp 打印帮助信息
func (c *Console) printHelp() {
	fmt.Println("Available commands:")
	fmt.Println("  help                - Show this help message")
	fmt.Println("  stats               - Show gateway statistics")
	fmt.Println("  broadcast <message> - Broadcast message to all online users")
	fmt.Println("  check <user_id>     - Check if user is online")
	fmt.Println("  exit/quit           - Exit console")
}

// printStats 打印网关统计信息
func (c *Console) printStats() {
	onlineCount := c.gateway.GetOnlineCount()
	fmt.Printf("Online users: %d\n", onlineCount)
}

// broadcastMessage 广播消息
func (c *Console) broadcastMessage(message string) {
	msg := types.Message{
		Header: types.MessageHeader{
			Type:      types.MessageTypeText,
			Timestamp: time.Now(),
			From:      "system",
		},
		Payload: []byte(message),
	}

	errs := c.gateway.Broadcast(msg)
	if len(errs) > 0 {
		fmt.Printf("Broadcast completed with %d errors\n", len(errs))
	} else {
		fmt.Println("Broadcast completed successfully")
	}
}

// checkUser 检查用户在线状态
func (c *Console) checkUser(userID string) {
	isOnline := c.gateway.IsUserOnline(userID)
	if isOnline {
		fmt.Printf("User %s is online\n", userID)
	} else {
		fmt.Printf("User %s is offline\n", userID)
	}
}
