package basetype

import (
	"encoding/binary"
	"fmt"
	"github.com/chyroc/go-redis/internal/helper"
	"github.com/chyroc/go-redis/internal/logger"
	"github.com/sirupsen/logrus"
)

// 整数集合

type Intset struct {
	encoding uint32
	length   uint32
	contents []byte
}

func (is *Intset) GetData(idx uint32) int64 {
	return getInt(&is.contents, idx, is.encoding)
}

func (is *Intset) Len() uint32 {
	return is.length
}

func (is *Intset) Int64Array() []int64 {
	var result []int64
	var i uint32
	for i = 0; i < is.length; i++ {
		result = append(result, getInt(&is.contents, i, is.encoding))
	}
	return result
}

const (
	INTSET_ENC_INT16 = 2
	INTSET_ENC_INT32 = 4
	INTSET_ENC_INT64 = 8
)

func NewIntset() *Intset {
	return &Intset{
		encoding: INTSET_ENC_INT16,
		length:   0,
		contents: make([]byte, 8),
	}
}

// s 可以是 int16,int32,int64，int 视为 int32
func (is *Intset) Add(s interface{}) {
	newEncoding := is.getDataEncoding(s)
	data := is.toInt64(s)

	// 扩容 & 升级
	is.expansion(newEncoding)

	// 找第一个大于等于 data 的 idx
	idx, found := helper.SearchFirstGreaterEqual(is, data)
	if !found {
		is.shiftBack(idx)
	}
	setInt(&is.contents, idx, newEncoding, data)
	if !found {
		is.length++
	}
	logrus.WithFields(logrus.Fields{
		"idx":    idx,
		"found":  found,
		"data":   s,
		"intset": is.Int64Array(),
	}).Info("[intset] set data end")
}

func (is *Intset) Exist(s int64) bool {
	var i uint32 = 0
	for i = 0; i < is.length; i++ {
		if d := getInt(&is.contents, i, is.encoding); d == s {
			return true
		}
	}
	return false
}

// 在 intset 中，求 小于等于 data 的第一个值的 idx
// 如果找到了，返回 idx, true
// 如果没有找到，返回小于 data 的第一个数据的 idx, false
func (is *Intset) find(data int64) (uint32, bool) {
	var left uint32 = 0
	var right = is.length - 1
	for {
		var mid = (left + right) / 2
		d := getInt(&is.contents, mid, is.encoding)
		if data == d {
			return mid, true
		}

		if data > d {
			left = mid
		} else if data < d {
			right = mid
		}

		if right-left == 1 {
			return left, false
		}
	}
}

// 将 idx 开始的数据后移一位，包括 idx 本身
func (is *Intset) shiftBack(idx uint32) {
	if idx >= is.length {
		return
	}
	copy(is.contents[(idx+1)*is.encoding:(is.length+1)*is.encoding], is.contents[idx*is.encoding:is.length*is.encoding])
}

func (is *Intset) getDataEncoding(i interface{}) uint32 {
	switch i.(type) {
	case int, int32:
		return INTSET_ENC_INT32
	case int16:
		return INTSET_ENC_INT16
	case int64:
		return INTSET_ENC_INT64
	}
	panic(fmt.Sprintf("%s(%T) 不允许添加到 intset 中", i, i))
}

func (is *Intset) toInt64(i interface{}) int64 {
	switch ii := i.(type) {
	case int:
		return int64(ii)
	case int32:
		return int64(ii)
	case int16:
		return int64(ii)
	case int64:
		return ii
	}
	panic(fmt.Sprintf("%s(%T) 不允许添加到 intset 中", i, i))
}

// 扩容 & 升级
// 扩容策略：
// - 首先计算需要的新的空间大小
// - 如果新空间大于(等于) 1024，那么升级为【大于新空间】的【1024 的整数倍】
// - 如果新空间小于 1024，那么升级为【新空间长度】的 【2 倍】
func (is *Intset) expansion(encoding uint32) {
	// 判断是否需要扩容和重新编码
	length := uint32(len(is.contents))
	newLength := encoding * (is.length + 1)
	if length < newLength {
		if newLength < 1024 {
			length = 2 * newLength
		} else {
			length = 1024 * (newLength/1024 + 1)
		}
	}
	var needExpansion = length > uint32(len(is.contents)) // 需要扩容
	var needUpgradeEncoding = encoding > is.encoding      // 每个数字的容量扩大了，需要重新编码，不能直接 copy bytes

	// 复制数据 & 重新编码
	if needUpgradeEncoding {
		var bs = make([]byte, length)
		var idx uint32 = 0
		for idx = 0; idx < is.length; idx++ {
			setInt(&bs, idx, encoding, getInt(&is.contents, idx, is.encoding))
		}
		is.contents = bs
		is.encoding = encoding
		return // 需要重新编码，此时不判断需不需要扩容，因为都需要制造新的 byte
	} else {
		if needExpansion {
			var bs = make([]byte, length)
			copy(bs, is.contents)
			is.contents = bs
			return // 需要扩容，不需要编码【直接复制 bytes】
		} else {
			return // 不需要扩容，不需要重新编码【直接返回】
		}
	}
}

func setInt(b *[]byte, idx uint32, encoding uint32, i int64) {
	logger.Log.WithFields(logrus.Fields{
		"idx":      idx,
		"encoding": encoding,
		"data":     i,
	}).Infof("[intset][set] set data")
	switch encoding {
	case 2:
		binary.LittleEndian.PutUint16((*b)[idx*encoding:(idx+1)*encoding], uint16(i))
	case 4:
		binary.LittleEndian.PutUint32((*b)[idx*encoding:(idx+1)*encoding], uint32(i))
	default:
		binary.LittleEndian.PutUint64((*b)[idx*encoding:(idx+1)*encoding], uint64(i))
	}
}

func getInt(b *[]byte, idx uint32, encoding uint32) int64 {
	logger.Log.WithFields(logrus.Fields{
		"idx":      idx,
		"encoding": encoding,
	}).Infof("[intset][get] get data")
	switch encoding {
	case 2:
		return int64(binary.LittleEndian.Uint16((*b)[idx*encoding : (idx+1)*encoding]))
	case 4:
		return int64(binary.LittleEndian.Uint32((*b)[idx*encoding : (idx+1)*encoding]))
	default:
		return int64(binary.LittleEndian.Uint64((*b)[idx*encoding : (idx+1)*encoding]))
	}
}
