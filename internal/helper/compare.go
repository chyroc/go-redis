package helper

import (
	"fmt"
	"reflect"
)

type CompareResult int

func (c CompareResult) IsEqual() bool {
	return c == 0
}

func (c CompareResult) IsBigger() bool {
	return c == 1
}

func (c CompareResult) IsSmall() bool {
	return c == -1
}

// a > b,  1
// a == b, 0
// a < b, -1
func Compare(a, b interface{}) CompareResult {
	at := reflect.TypeOf(a)
	bt := reflect.TypeOf(b)
	if at.Kind() != bt.Kind() {
		panic(fmt.Sprintf("参与比较的两个参数的类型不一致(%s != %s)", at.Kind(), bt.Kind()))
	}

	switch at.Kind() {
	case reflect.Uint:
		return tocompare(a.(uint) > b.(uint), b.(uint) > a.(uint))
	case reflect.Uint8:
		return tocompare(a.(uint8) > b.(uint8), b.(uint8) > a.(uint8))
	case reflect.Uint16:
		return tocompare(a.(uint16) > b.(uint16), b.(uint16) > a.(uint16))
	case reflect.Uint32:
		return tocompare(a.(uint32) > b.(uint32), b.(uint32) > a.(uint32))
	case reflect.Uint64:
		return tocompare(a.(uint64) > b.(uint64), b.(uint64) > a.(uint64))
	case reflect.Int:
		return tocompare(a.(int) > b.(int), b.(int) > a.(int))
	case reflect.Int8:
		return tocompare(a.(int8) > b.(int8), b.(int8) > a.(int8))
	case reflect.Int16:
		return tocompare(a.(int16) > b.(int16), b.(int16) > a.(int16))
	case reflect.Int32:
		return tocompare(a.(int32) > b.(int32), b.(int32) > a.(int32))
	case reflect.Int64:
		return tocompare(a.(int64) > b.(int64), b.(int64) > a.(int64))
	case reflect.Float32:
		return tocompare(a.(float32) > b.(float32), b.(float32) > a.(float32))
	case reflect.Float64:
		return tocompare(a.(float64) > b.(float64), b.(float64) > a.(float64))
	case reflect.Bool:
		return tocompare(a.(bool) && !b.(bool), b.(bool) && !a.(bool))
	case reflect.String:
		for k := 0; k < len(a.(string)) && k < len(b.(string)); k++ {
			if d := Compare(a.(string)[k], b.(string)[k]); d != 0 {
				return d
			}
		}
		return Compare(len(a.(string)), len(b.(string)))
	}
	panic(fmt.Sprintf("%s 是不支持比较的类型", at.Kind()))
}

func tocompare(aBiggerThenB bool, bBiggerThenA bool) CompareResult {
	if aBiggerThenB {
		return 1
	}
	if bBiggerThenA {
		return -1
	}
	return 0
}
