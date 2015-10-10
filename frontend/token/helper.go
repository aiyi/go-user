package token

import (
	"github.com/chanxuehong/util/random"
)

func NewTokenId() string {
	return string(random.NewRandomEx())
}

func ExpirationAccess(timestamp int64) int64 {
	return timestamp + 7200
}

func ExpirationRefresh(timestamp int64) int64 {
	return timestamp + 31556952
}
