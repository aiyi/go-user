package mc

import (
	"github.com/aiyi/go-user/config"
)

func getServerList() ([]string, error) {
	return config.ConfigData.MemcacheServerList, nil
}
