package utils

import (
	"github.com/bwmarrin/snowflake"
	"sync"
	"time"
)

var (
	node *snowflake.Node
	once sync.Once
	err  error
)

func GenerateSnowflakeID() int64 {
	once.Do(func() {
		node, err = snowflake.NewNode(1)
		if err != nil {
			return
		}
	})

	if err != nil {
		return time.Now().UnixNano()
	}

	id := node.Generate()
	return id.Int64()
}
