package internal

type zskiplistNode struct {
	level []struct {
		forward *zskiplistNode
		span    uint32
	}
	backward *zskiplistNode
	score    float64
	obj      interface{}
}

type zskiplist struct {
	head   *zskiplistNode
	tail   *zskiplistNode
	length uint32
	level  int
}
