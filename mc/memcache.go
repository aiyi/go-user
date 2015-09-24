package mc

import (
	"github.com/bradfitz/gomemcache/memcache"
)

var client *memcache.Client

func init() {
	servers, err := getServerList()
	if err != nil {
		panic(err)
	}

	client = memcache.New(servers...)
}

func Client() *memcache.Client {
	return client
}
