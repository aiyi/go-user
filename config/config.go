package config

import (
	"encoding/json"
	"io/ioutil"
)

func init() {
	if err := loadConfig(&ConfigData); err != nil {
		panic(err)
	}
}

var ConfigData Config // read only

type Config struct {
	Mysql struct {
		UserName string `json:"username"`
		Password string `json:"password"`
		Host     string `json:"host"`
		Port     int    `json:"port"`
		Database string `json:"database"`
	} `json:"mysql"`
	MemcacheServerList []string `json:"memcache_server_list,omitempty"`
	SecurityKey        []byte   `json:"security_key,omitempty"`
}

func loadConfig(cfg *Config) (err error) {
	configBytes, err := ioutil.ReadFile("config/config.json.bak")
	if err != nil {
		return
	}
	return json.Unmarshal(configBytes, cfg)
}
