package mc

import (
	"github.com/bradfitz/gomemcache/memcache"
)

var client = memcache.New("xxx.xxx.xxx.xxx:11211")

func Client() *memcache.Client {
	return client
}
