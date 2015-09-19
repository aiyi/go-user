// 提供生成 userid 的 api.
//
// 需要注意 SnowflakeWorkerId 不能重复!
package userid

import (
	"github.com/chanxuehong/util/id"
)

func init() {
	snowflakeWorkerId, err := getSnowflakeWorkerId()
	if err != nil {
		panic(err)
	}

	err = id.SetSnowflakeWorkerId(snowflakeWorkerId)
	if err != nil {
		panic(err)
	}
}

// 集群环境下不能重复
func getSnowflakeWorkerId() (int, error) {
	return 0, nil
}

func GetId() (userid int64, err error) {
	return id.NewSnowflakeId()
}
