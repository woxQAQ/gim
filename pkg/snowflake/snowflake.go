package snowflake

import (
	"sync"
	"time"

	"github.com/bwmarrin/snowflake"
	"github.com/spf13/viper"
	"github.com/woxQAQ/gim/pkg/constants"
)

var (
	node *snowflake.Node
	once sync.Once
)

// GenerateID 生成唯一标识
func GenerateID() string {
	once.Do(func() {
		var err error
		// 如果未初始化，使用默认节点
		nodeId := viper.GetInt64(constants.SNOWFLAKE_NODE_ID)
		node, err = snowflake.NewNode(nodeId)
		if err != nil {
			panic(err)
		}
	})
	if node == nil {
		// 如果节点未初始化，使用时间戳作为备选方案
		return time.Now().Format("20060102150405") + "-" + time.Now().Format("000000")
	}
	return node.Generate().String()
}
