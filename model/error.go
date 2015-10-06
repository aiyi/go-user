package model

import (
	"errors"
)

var ErrNotFound = errors.New("item not found") // 所有的没有找到的都返回这个错误
