package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// statsCmd represents the stats command
var statsCmd = &cobra.Command{
	Use:   "stats [user_id] [platform_id]",
	Short: "显示网关统计信息",
	Long: `stats 命令用于显示WebSocket网关的统计信息。

用法：
  stats                     显示总体统计信息
  stats <user_id>           显示指定用户的统计信息
  stats <user_id> <platform_id>  显示指定用户在特定平台的统计信息`,
	Run: func(cmd *cobra.Command, args []string) {
		// 显示总体统计信息
		onlineCount := gateway.GetOnlineCount()
		fmt.Printf("在线用户数: %d\n", onlineCount)

		// 如果提供了用户ID，显示用户特定信息
		if len(args) > 0 {
			userID := args[0]
			platformID := int32(1) // 默认平台ID
			if len(args) > 1 {
				// 尝试解析平台ID
				if id, err := fmt.Sscanf(args[1], "%d", &platformID); err != nil || id != 1 {
					fmt.Printf("无效的平台ID: %s\n", args[1])
					return
				}
			}

			// 获取用户在线状态
			isOnline := gateway.IsUserOnline(userID)
			fmt.Printf("用户 %s 状态: %s\n", userID, map[bool]string{true: "在线", false: "离线"}[isOnline])

			// 如果用户在线，获取心跳状态
			if isOnline {
				lastPing, err := gateway.GetUserHeartbeatStatus(userID, platformID)
				if err != nil {
					fmt.Printf("获取用户心跳状态失败: %v\n", err)
				} else {
					fmt.Printf("最后心跳时间: %s\n", lastPing.Format("2006-01-02 15:04:05"))
				}
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(statsCmd)
}
