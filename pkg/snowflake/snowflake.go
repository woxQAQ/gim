package snowflake

import (
	"time"

	"github.com/bwmarrin/snowflake"
)

var node *snowflake.Node

// InitGenerator 初始化ID生成器
func InitGenerator(nodeID int64) error {
	var err error
	node, err = snowflake.NewNode(nodeID)
	return err
}

// GenerateID 生成唯一标识
func GenerateID() string {
	if node == nil {
		// 如果节点未初始化，使用时间戳作为备选方案
		return time.Now().Format("20060102150405") + "-" + time.Now().Format("000000")
	}
	return node.Generate().String()
}
