// 提供生成 userid 的 api.
package userid

import (
	"github.com/chanxuehong/util/id"
)

func init() {
	snowflakeWorkerId, err := getSnowflakeWorkerId()
	if err != nil {
		panic(err)
	}

	if err = id.SetSnowflakeWorkerId(snowflakeWorkerId); err != nil {
		panic(err)
	}
}

func GetId() (userid int64, err error) {
	return id.NewSnowflakeId()
}
