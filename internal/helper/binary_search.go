package helper

import (
	"github.com/chyroc/go-redis/internal/logger"
	"github.com/sirupsen/logrus"
)

// 顺序的序列
// - 可以获取长度
// - 可以根据偏移获取指定的值
type SequentialSequence interface {
	GetData(idx uint32) int64
	Len() uint32
}

// 搜索顺序序列中，第一个大于等于 data 的序号，返回是 [0, n]
// - 为空，返回 0
// - 比所有值都小，返回 0
// - 比所有值都大，返回 n
// - 至少比一个值大，但是又不是比所有的值都大，返回范围是 [1, n-1]
func SearchFirstGreaterEqual(ss SequentialSequence, data int64) (idx uint32, found bool) {
	n := ss.Len()
	if n == 0 {
		return 0, false
	} else if data < ss.GetData(0) {
		return 0, false
	} else if data > ss.GetData(n-1) {
		return n, false
	}

	var low uint32 = 0
	var high = n - 1
	for low <= high {
		// 避免溢出
		mid := low + (high-low)>>1 // 0~x 的范围
		d := ss.GetData(mid)
		logger.Log.WithFields(logrus.Fields{
			"low":    low,
			"high":   high,
			"mid":    mid,
			"search": data,
			"got":   d,
		}).Infof("[helper.binary_search] get mid data")
		if d >= data {
			if mid == 0 || ss.GetData(mid-1) < data {
				return mid, data == d
			} else {
				high = mid - 1
			}
		} else {
			low = mid + 1
		}
	}

	panic("不可能走到这里")
}

type IntSequentialSequence []int

func (r IntSequentialSequence) GetData(idx uint32) int64 {
	return int64(r[idx])
}

func (r IntSequentialSequence) Len() uint32 {
	return uint32(len(r))
}
