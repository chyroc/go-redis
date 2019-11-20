package basetype

import (
	"fmt"
	"github.com/chyroc/go-redis/internal/helper"
)

// todo 暂时固定为 1
const rehashLoadFactor = 1

type EntrySetter interface {
	SetEntry(entry *DictEntry)
}

type DictEntry struct {
	key   *SDS
	value interface{}
	next  *DictEntry
}

func (d *DictEntry) SetEntry(entry *DictEntry) {
	if d.next != nil && d.next.next != nil {
		if entry == nil {
			d.next = d.next.next
			return
		}
		entry.next = d.next.next
	}
	d.next = entry
}

type dhtEntrySeter struct {
	dht   *DictHashTable
	index uint64
}

func (ds *dhtEntrySeter) SetEntry(entry *DictEntry) {
	e := ds.dht.table[ds.index]
	if e == nil || e.next == nil {
		ds.dht.table[ds.index] = entry
		return
	}
	if entry != nil {
		entry.next = e.next.next
		ds.dht.table[ds.index] = entry
		return
	}
	ds.dht.table[ds.index] = ds.dht.table[ds.index].next
}

type DictHashTable struct {
	table    []*DictEntry
	size     uint32
	sizemask uint32
	used     uint32
}

func (dht *DictHashTable) keyIndex(k string) uint64 {
	khash := helper.Hash([]byte(k))
	return uint64(dht.sizemask) & khash
}

func (dht *DictHashTable) loadFactor() float64 {
	return float64(dht.used) / float64(dht.size)
}

type DictType struct {
}

type Dict struct {
	type_    *DictType
	privdata interface{}
	ht       [2]*DictHashTable

	// 没有 rehash 的时候为 -1，开始 rehash 的时候为 0，每进行一次 key 的操作，就 +1
	reHashIndex int
}

func newHT(initSize uint32) *DictHashTable {
	return &DictHashTable{
		table:    make([]*DictEntry, initSize),
		size:     initSize,
		sizemask: initSize - 1,
		used:     0,
	}
}

func NewDict() *Dict {
	d := &Dict{
		type_:       nil,
		privdata:    nil,
		ht:          [2]*DictHashTable{newHT(8), nil},
		reHashIndex: -1,
	}
	return d
}

func (d *Dict) Set(k string, v interface{}) {
	if !d.isRehashing() && d.ht[0].loadFactor() > rehashLoadFactor {
		d.startRehash()
	}
	if d.isRehashing() && d.ht[0].used == 0 {
		d.stopRehash()
	}

	// value
	nextEntry := &DictEntry{
		key:   NewSDSWithString(k),
		value: v,
		next:  nil,
	}

	if !d.isRehashing() {
		entry, entrySetter := d.findEntrySetter(0, k)
		if entry == nil {
			d.ht[0].used++
		}
		entrySetter.SetEntry(nextEntry)
	} else {
		// 开始 rehash 了
		if entry, entrySetter := d.findEntrySetter(0, k); entry != nil {
			d.reHashIndex++
			d.ht[0].used--
			entrySetter.SetEntry(nil)
		}

		entry, entrySetter := d.findEntrySetter(1, k)
		if entry == nil {
			d.ht[1].used++
		}
		entrySetter.SetEntry(nextEntry)
	}
}

func (d *Dict) Size() uint32 {
	if d.isRehashing() {
		return d.ht[0].used + d.ht[1].used
	} else {
		return d.ht[0].used
	}
}

// todo: rehash
func (d *Dict) Get(k string) interface{} {
	if d.isRehashing() && d.ht[0].used == 0 {
		d.stopRehash()
	}

	if !d.isRehashing() {
		if entry, _ := d.findEntrySetter(0, k); entry != nil {
			return entry.value
		}
	} else {
		// 开始 rehash 了

		if entry, _ := d.findEntrySetter(0, k); entry != nil {
			d.Set(k, entry.value) // set 会 rehash 这个 key
			return entry.value
		}

		if entry, _ := d.findEntrySetter(1, k); entry != nil {
			return entry.value
		}
	}

	return nil
}

func (d *Dict) Del(k string) {
	if entry, entrySetter := d.findEntrySetter(0, k); entry != nil {
		d.ht[0].used--
		entrySetter.SetEntry(nil)
	}

	if d.isRehashing() {
		if entry, entrySetter := d.findEntrySetter(1, k); entry != nil {
			d.ht[1].used--
			entrySetter.SetEntry(nil)
		}
	}
}

func (d *Dict) findEntrySetter(tableIndex int, k string) (entry *DictEntry, entrySetter EntrySetter) {
	kindex := d.ht[tableIndex].keyIndex(k)
	entry = d.ht[tableIndex].table[kindex]
	if entry == nil {
		return nil, &dhtEntrySeter{
			dht:   d.ht[tableIndex],
			index: kindex,
		}
	}

	var pre *DictEntry
	for entry != nil {
		if entry.key.EqualToString(k) {
			if pre == nil {
				return entry, &dhtEntrySeter{
					dht:   d.ht[tableIndex],
					index: kindex,
				}
			}
			return entry, pre
		}
		pre = entry
		entry = pre.next
	}
	return nil, pre
}

func (d *Dict) isRehashing() bool {
	return d.reHashIndex != -1
}

func (d *Dict) startRehash() {
	d.reHashIndex = 0 // 开始 rehash
	if d.ht[0].used < 1024 {
		d.ht[1] = newHT(d.ht[0].used * 2)
	} else {
		d.ht[1] = newHT(d.ht[0].used + 1024)
	}

	fmt.Printf("[===] 开始 rehash\n")
}

func (d *Dict) stopRehash() {
	d.reHashIndex = -1
	d.ht[0] = d.ht[1]
	d.ht[1] = nil

	fmt.Printf("[===] 停止 rehash\n")
}
