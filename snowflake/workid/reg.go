/**
 * @Author:huangshan
 * @Description: 注册机制获取机器ID
 * @Date: 2022/11/3 16:30
 */
package workid

import (
	"fmt"
	"os"
	"strconv"
)

// 通过注册获取机器ID
type Reg struct {
	workId int64 //环境变量中的workid
}

const (
	WORK_ENV = "work_id"
)

func (r *Reg) GetWorkId() (int64, error) {
	//说明第一次获取，从环境变量中拿
	if r.workId <= 0 {
		str := os.Getenv(WORK_ENV)
		if str == "" {
			//return 0, fmt.Errorf("env (%s) val empty", WORK_ENV)
			str = "1"
		}
		i, err := strconv.Atoi(str)
		if err != nil {
			fmt.Printf("[Reg][GetWorkId] get work id failed ,err(%+v)", err)
			return 0, err
		}
		r.workId = int64(i)
		fmt.Println("环境变量获取 ", r.workId)
		return r.workId, nil
	}

	//调用api主动获取
	re := &Random{}
	workId, _ := re.GetWorkId()
	fmt.Println("api获取 ", workId)
	return workId, nil
}
