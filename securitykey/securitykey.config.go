package securitykey

import (
	"github.com/aiyi/go-user/config"
)

// 从私密存储上获取安全key
func getKey() ([]byte, error) {
	return config.ConfigData.SecurityKey, nil
}
