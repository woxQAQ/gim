package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// statsCmd represents the stats command
var statsCmd = &cobra.Command{
	Use:   "stats",
	Short: "显示网关统计信息",
	Long: `stats 命令用于显示WebSocket网关的统计信息，包括：
- 当前在线用户数量
- 其他统计指标`,
	Run: func(cmd *cobra.Command, args []string) {
		onlineCount := gateway.GetOnlineCount()
		fmt.Printf("在线用户数: %d\n", onlineCount)
	},
}

func init() {
	rootCmd.AddCommand(statsCmd)
}
