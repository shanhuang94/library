/**
 * @Author:huangshan
 * @Description: 生成雪花算法
 * @Date: 2022/11/2 11:51
 */
package snowflake

import (
	"fmt"
	"math"
	"neuroxess-cloud/library/snowflake/workid"
	"sync"
	"time"
)

const (
	epoch          = int64(1667457275000)              // 设置起始时间(时间戳/毫秒)：2020-01-01 00:00:00，有效期69年
	timestampBits  = uint(41)                          // 时间戳占用位数
	workeridBits   = uint(10)                          // 机器id所占位数
	sequenceBits   = uint(12)                          // 序列所占的位数
	timestampMax   = int64(-1 ^ (-1 << timestampBits)) // 时间戳最大值
	workeridMax    = int64(-1 ^ (-1 << workeridBits))  // 支持的最大机器id数量
	sequenceMask   = int64(-1 ^ (-1 << sequenceBits))  // 支持的最大序列id数量
	workeridShift  = sequenceBits                      // 机器id左移位数
	timestampShift = sequenceBits + workeridBits       // 时间戳左移位数
	delayThreshold = 60                                //延迟阈值 60毫秒
)

var s *snowflake
var w workid.WorkId

type snowflake struct {
	sync.Mutex          // 锁
	lastTimestamp int64 // 上次时间戳 ，毫秒
	workerId      int64 // 工作节点
	sequence      int64 // 序列号
}

func Init(work workid.WorkId) error {
	//获取机器ID
	w = work
	id, err := w.GetWorkId()
	if err != nil {
		return err
	}
	s = &snowflake{
		lastTimestamp: 0,
		sequence:      0,
		workerId:      id,
	}
	return nil
}

func GetWorkId(work workid.WorkId) int64 {
	res, err := work.GetWorkId()
	if err != nil {
		return -1
	}
	return res
}

func NextVal() int64 {
	s.Lock()
	var err error
	now := timeGen()
	//极端case 时间回拨判断
	if now < s.lastTimestamp {
		diff := s.lastTimestamp - now
		//情况1 如果回拨时间在阈值内，等待时钟对齐
		if diff < delayThreshold {
			waitClkAlign(s.lastTimestamp)
		} else {
			//情况2 时间太长，直接换一个机器ID
			s.workerId, err = w.GetWorkId()
			if err != nil {
				s.Unlock()
				//todo 告警，需要运维开发知晓
				return 0
			}
		}
		now = timeGen()
	}

	//同一毫秒
	if s.lastTimestamp == now {
		// 当同一时间戳（精度：毫秒）下多次生成id会增加序列号
		s.sequence = (s.sequence + 1) & sequenceMask
		// 如果当前序列超出12bit长度，则需要等待下一毫秒
		if s.sequence == 0 {
			// 等一毫秒 下一毫秒将使用sequence:0
			for now == s.lastTimestamp {
				now = timeGen()
			}
		}
	} else {
		// 不同时间戳（精度：毫秒）下直接使用序列号：0
		s.sequence = 0
	}
	t := now - epoch
	if t > timestampMax {
		s.Unlock()
		//todo 告警
		return 0
	}
	s.lastTimestamp = now
	r := t<<timestampShift | (s.workerId << workeridShift) | (s.sequence)
	s.Unlock()
	return r
}

func timeGen() int64 {
	return time.Now().UnixNano() / 1000000 // 转毫秒
}

// 等待时钟对齐
func waitClkAlign(lastTimestamp int64) {
	now := timeGen()
	for lastTimestamp > now {
		time.Sleep(time.Millisecond)
		now = timeGen()
	}
}

// 计算 不同bit分配 时间和序列情况
func turn(time int, danwei string) string {
	t := math.Pow(float64(2), float64(time))
	s := 44 - time
	se := math.Pow(float64(2), float64(s))
	tmp := float64(0)
	switch danwei {
	case "毫秒":
		tmp = 1000
	case "秒":
		tmp = 1
	case "10毫秒":
		tmp = 100
	case "100毫秒":
		tmp = 10
	}
	y := float64(t) / 3600 / 24 / 365 / tmp
	return fmt.Sprintf("时间是%d位，单位是%s，可以用%f年，序列号可以有%f个", time, danwei, y, se)
}
