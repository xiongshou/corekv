package utils

import (
	"bytes"
	"core/utils/codec"
	"math/rand"
	"sync"
)

const (
	defaultMaxHeight = 48
)

type Element struct{
	entry *codec.Entry
	levels []*Element
	score float64
}

type SkipList struct {
	header *Element

	rand *rand.Rand
	maxLevel int
	lock sync.RWMutex
	size int64
	length int
}

func newElement(score float64, entry *codec.Entry, level int) *Element{
	return &Element{
		levels: make([]*Element, level + 1),
		entry: entry,
		score: score,
	}
}
func (elem *Element) Entry() *codec.Entry {
	return elem.entry
}

func NewSkipList() *SkipList{
	header := &Element{
		levels: make([]*Element, defaultMaxHeight),
	}
	return &SkipList{
		header: header,
		maxLevel: defaultMaxHeight - 1,
		rand: r,
	}
}

func (list *SkipList) Add(data *codec.Entry) error{
	list.lock.Lock()
	defer list.lock.Unlock()

	prevs := make([]*Element, list.maxLevel + 1)

	key := data.Key
	keyScore := list.calcScore(key)
	header, maxLevel := list.header, list.maxLevel
	prev := header

	for i := maxLevel; i>=0; i--{
		for ne := prev.levels[i]; ne != nil; ne = prev.levels[i] {
			if comp := list.compare(keyScore, key, ne); comp <= 0{
				if comp == 0{
					ne.entry = data
					return nil
				} else{
					//  ne节点小
					prev = ne
				}
			}else{
				break
			}
		}
		// prev - next - 最后
		//  prevs: key 的前一个节点
		prevs[i] = prev
	}

	randLevel, keyScore := list.randLevel(), list.calcScore(key)
	e := newElement(keyScore, data, randLevel)

	// preves : 每层第一个大于等于 key 的element
	for i:= randLevel; i>=0; i--{
		ne := prevs[i].levels[i]
		prevs[i].levels[i] = e
		e.levels[i] = ne
	}

	return nil
}

func (list *SkipList) Search(key []byte) (e *codec.Entry) {

	list.lock.RLock()
	defer list.lock.RUnlock()
	keyScore := list.calcScore(key)

	header, maxLevel := list.header, list.maxLevel

	prev := header
	for i:= maxLevel; i>=0 ;i-- {
		for ne := prev.levels[i]; ne!=nil; ne= prev.levels[i] {
			if comp := list.compare(keyScore, key, ne); comp <= 0 {
				if comp == 0 {
					return ne.entry
				} else {
					prev = ne
				}
			}else{
				break
			}
		}
	}
	return nil

}
func (list *SkipList) calcScore(key []byte)(score float64){
	var hash uint64
	l := len(key)

	if(l>8){
		l=8
	}

	for i:=0; i<l; i++ {
		shift := uint(64 - 8 - i*8)
		hash |= uint64(key[i]) << shift
	}
	score = float64(hash)
	return
}



func (list *SkipList)compare(score float64, key []byte, next * Element) int{
	if score == next.score {
		return bytes.Compare(key, next.entry.Key)
	}

	if score < next.score {
		return -1
	} else {
		return 1
	}
}

func (list *SkipList) Size() int64{
	return list.size
}

func (list *SkipList) randLevel() int{
	for i:=0; i< list.maxLevel; i++ {
		if list.rand.Intn(2) == 0 {
			return i
		}
	}
	return list.maxLevel
}