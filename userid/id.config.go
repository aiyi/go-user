package userid

import (
	"github.com/aiyi/go-user/config"
)

// 集群环境下不能重复
func getSnowflakeWorkerId() (int, error) {
	return config.ConfigData.SnowflakeWorkerId, nil
}
