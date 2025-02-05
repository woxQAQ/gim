package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// checkCmd represents the check command
var checkCmd = &cobra.Command{
	Use:   "check [user_id]",
	Short: "检查用户在线状态",
	Long: `check 命令用于检查指定用户的在线状态。

示例：
  console check user123
  console check admin`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		userID := args[0]
		isOnline := gateway.IsUserOnline(userID)
		if isOnline {
			fmt.Printf("用户 %s 当前在线\n", userID)
		} else {
			fmt.Printf("用户 %s 当前离线\n", userID)
		}
	},
}

func init() {
	rootCmd.AddCommand(checkCmd)
}
