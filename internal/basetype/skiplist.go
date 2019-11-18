package basetype

import (
	"errors"
	"fmt"
	"github.com/chyroc/go-redis/internal/helper"
	"github.com/chyroc/go-redis/internal/logger"
	"github.com/sirupsen/logrus"
	mathrand "math/rand"
	"strings"
	"time"
)

type SkipListNodeLevel struct {
	forward *SkipListNode
	span    uint32
}

type SkipListNode struct {
	level    []*SkipListNodeLevel
	backward *SkipListNode
	score    float64
	obj      string
}

// 先以 score 排序找，然后以 key 排序找
func (node *SkipListNode) isTail() bool {
	return len(node.level) == 32 && node.level[0] == nil
}

// 先以 score 排序找，然后以 key 排序找
func (node *SkipListNode) getNodeOrNext(key string, score float64) *SkipListNode {
	return node.getNodeOrNextWithCompare(key, score)
}

func (node *SkipListNode) getNodeOrNextWithCompare(key string, score float64) *SkipListNode {
	// 0个node，返回 tail，true
	// 有这个 score，返回这个node，true
	// 比某些大，比某些小，返回大于这个score的第一个 node，true
	// 比所有数据都大，返回tail，true
	// 比所有数据都小，返回第一个 node，true（这个实际就是大于这个score的第一个node，归类到第三种）

	var compare = func(s string, f float64, s2 string, f2 float64) helper.CompareResult {
		if result := helper.Compare(f, f2); !result.IsEqual() {
			return result
		}
		return helper.Compare(s, s2)
	}

	for i := len(node.level) - 1; i >= 0; i-- {
		//logger.Log.WithFields(logrus.Fields{"i": i, "level[i]": node.level[i],}).Infof("[skiplist][get] getNodeOrNextWithCompare")

		if node.level[i] == nil || node.level[i].forward == nil {
			if i == 0 {
				return nil // 这里应该是整个 skiplist 的最后一个 node 节点，只能在 skiplist 上获取，这里用 nil 表示
			}
			continue
		}

		next := node.level[i].forward
		if next.isTail() {
			if i == 0 {
				return nil // 这里应该是整个 skiplist 的最后一个 node 节点，只能在 skiplist 上获取，这里用 nil 表示
			}
			continue
		}

		scoreCompare := compare(key, score, next.obj, next.score)
		logger.Log.WithField("scoreCompare", scoreCompare).Infof("[skiplist][get] getNodeOrNextWithCompare")

		if scoreCompare.IsEqual() {
			return next
		} else if scoreCompare.IsBigger() {
			logger.Log.WithField("next", next).WithField("next.score", next.score).Infof("[skiplist][get] found next, goto next")
			return next.getNodeOrNextWithCompare(key, score)
		}

		// next 不为空；score 小于 next；在最底层链表
		if i == 0 {
			return next
		}
	}

	panic("不可能走到这里")
}

type SkipList struct {
	head   *SkipListNode
	tail   *SkipListNode
	length uint32
	level  int
}

func (s *SkipList) String() string {
	b := new(strings.Builder)
	b.WriteString("SkipList\n")

	for i := 31; i >= 0; i-- {
		var node = s.head.level[i].forward
		j := 0
		add := false
		for {
			if node == nil {
				break
			}
			if node.isTail() {
				break
			}
			add = true
			if j == 0 {
				b.WriteString(fmt.Sprintf("%s:%s ", helper.String(i), helper.String(j)))
			}
			if j != 0 {
				b.WriteString(" -> ")
			}
			b.WriteString(fmt.Sprintf("[%d]<%s: %q>", len(node.level), helper.String(node.score), node.obj))
			if node.level[i] == nil {
				break
			}
			node = node.level[i].forward
			j++
		}
		if add {
			b.WriteString("\n")
		}
	}
	return b.String()
}

func NewSkipList() *SkipList {
	head := &SkipListNode{
		level:    make([]*SkipListNodeLevel, 32),
		backward: nil,
		score:    0,
		obj:      "",
	}
	tail := &SkipListNode{
		level:    make([]*SkipListNodeLevel, 32),
		backward: nil,
		score:    0,
		obj:      "",
	}
	tail.backward = head
	for i := 0; i < 32; i++ {
		head.level[i] = &SkipListNodeLevel{
			forward: tail,
			span:    0,
		}
	}

	return &SkipList{
		head: head,
		tail: tail,
	}
}

var ErrDataRepeated = errors.New("skiplist data repeated error")

func (s *SkipList) Add(key string, score float64) error {
	// 找到这个节点的下个节点
	nodeOrNext := s.find(key, score)
	logger.Log.WithFields(logrus.Fields{
		"key":              key,
		"score":            score,
		"nodeOrNext.score": nodeOrNext.score,
		"nodeOrNext.obj":   nodeOrNext.obj,
		"isTail":           nodeOrNext.isTail(),
	}).Infof("[skiplist][add] found node")
	if nodeOrNext.score == score && nodeOrNext.obj == key {
		return ErrDataRepeated
	}

	// 随机生成一个高度，[1-32] 之间
	level := genLevelHigh()
	logger.Log.WithFields(logrus.Fields{
		"key":   key,
		"score": score,
		"level": level,
	}).Infof("[skiplist][add] gen level")

	// 新节点
	newNode := &SkipListNode{
		level:    make([]*SkipListNodeLevel, level),
		backward: nil,
		score:    score,
		obj:      key,
	}

	// 获取前后节点
	//nextNode := nodeOrNext
	prevNode := nodeOrNext.backward

	// 4 个情况
	// 0,1,2
	// 0,2,1
	// 2,1,0
	// 2,0,1

	// 求 [0, N) 这 N 个 level 的前后节点
	// level -> [prev, next]
	nodes := make(map[int][2]*SkipListNode)
	for i := 0; i < level; i++ {
		for i >= len(prevNode.level) {
			prevNode = prevNode.backward
		}

		nodes[i] = [2]*SkipListNode{prevNode, prevNode.level[i].forward}
	}

	// 遍历 nodes，设置指针
	newNode.backward = nodes[0][0]
	for i := 0; i < level; i++ {
		nodes[i][0].level[i].forward = newNode
		newNode.level[i] = &SkipListNodeLevel{
			forward: nodes[i][1],
			span:    0,
		}
		nodes[i][1].backward = newNode
	}

	return nil
}

func (s *SkipList) Get(score float64) (result []*SkipListNode, err error) {
	logger.Log.WithFields(logrus.Fields{"score": score,}).Infof("[skiplist] get")

	nodeOrNext := s.head.getNodeOrNextWithCompare("", score)
	logger.Log.WithFields(logrus.Fields{"nodeOrNext": nodeOrNext,}).Infof("[skiplist] get")
	if nodeOrNext == nil {
		return
	}
	for nodeOrNext.score == score {
		result = append(result, nodeOrNext)
		nodeOrNext = nodeOrNext.level[0].forward
	}
	return
}

// preNode, currentNode, found
func (s *SkipList) find(key string, score float64) *SkipListNode {
	nodeOrNext := s.head.getNodeOrNext(key, score)
	if nodeOrNext == nil {
		nodeOrNext = s.tail
	}
	return nodeOrNext
}

func genLevelHigh() int {
	level := 1
	for level < 32 && mathrand.Intn(100) < 25 {
		level++
	}
	return level
}

func init() {
	mathrand.Seed(time.Now().UnixNano())
}
