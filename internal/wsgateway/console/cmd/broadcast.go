package cmd

import (
	"fmt"
	"strings"
	"time"

	"github.com/spf13/cobra"

	"github.com/woxQAQ/gim/internal/types"
)

// broadcastCmd represents the broadcast command
var broadcastCmd = &cobra.Command{
	Use:   "broadcast [message]",
	Short: "向所有在线用户广播消息",
	Long: `broadcast 命令用于向所有在线用户广播消息。

示例：
  console broadcast "Hello, everyone!"
  console broadcast "System maintenance in 5 minutes"`,
	Args: cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		message := strings.Join(args, " ")
		msg := types.Message{
			Header: types.MessageHeader{
				Type:      types.MessageTypeText,
				Timestamp: time.Now(),
				From:      "system",
			},
			Payload: []byte(message),
		}

		errs := gateway.Broadcast(msg)
		if len(errs) > 0 {
			fmt.Printf("广播完成，但有 %d 个错误发生\n", len(errs))
		} else {
			fmt.Println("广播成功完成")
		}
	},
}

func init() {
	rootCmd.AddCommand(broadcastCmd)
}
