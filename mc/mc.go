package mc

import (
	"github.com/bradfitz/gomemcache/memcache"
)

func init() {
	servers, err := getServerList()
	if err != nil {
		panic(err)
	}

	client = memcache.New(servers...)
}

var client *memcache.Client

func Client() *memcache.Client {
	return client
}
