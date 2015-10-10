package session

import (
	"github.com/chanxuehong/util/id"
)

// ^[A-Za-z0-9_-]+$
func NewSessionId() (sid string, err error) {
	sidx, err := id.NewSessionId()
	if err != nil {
		return
	}
	sid = string(sidx)
	return
}

// ^temp\.[A-Za-z0-9_-]+$
func NewGuestSessionId() (sid string, err error) {
	sidx, err := id.NewSessionId()
	if err != nil {
		return
	}
	sid = "temp." + string(sidx)
	return
}
