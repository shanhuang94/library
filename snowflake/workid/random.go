/**
 * @Author:huangshan
 * @Description:随机数生成机器ID
 * @Date: 2022/11/3 16:22
 */
package workid

import (
	"math/rand"
	"time"
)

type Random struct {
}

func (r *Random) GetWorkId() (int64, error) {
	rand.Seed(time.Now().UnixNano())
	return rand.Int63n(100), nil
}
