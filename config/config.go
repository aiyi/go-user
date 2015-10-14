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
	WebServer struct {
		BaseURL string `json:"base_url"` // http://localhost:8080
	} `json:"web_server"`

	Mysql struct {
		UserName string `json:"username"`
		Password string `json:"password"`
		Host     string `json:"host"`
		Port     int    `json:"port"`
		Database string `json:"database"`
	} `json:"mysql"`

	Weixin struct {
		MP struct { // 公众号
			AppId     string `json:"appid"`
			AppSecret string `json:"appsecret"`
		} `json:"mp"`
		Open struct { // 开放平台
			Web struct { // 网站应用
				AppId     string `json:"appid"`
				AppSecret string `json:"appsecret"`
			} `json:"web"`
		} `json:"open"`
	} `json:"weixin"`

	MemcacheServerList []string `json:"memcache_server_list,omitempty"`
	SecurityKey        []byte   `json:"security_key,omitempty"`
	SnowflakeWorkerId  int      `json:"snowflake_workerid"` // SnowflakeWorkerId 不能重复!
}

func loadConfig(cfg *Config) (err error) {
	configBytes, err := ioutil.ReadFile("config/config.json.bak")
	if err != nil {
		return
	}
	return json.Unmarshal(configBytes, cfg)
}
