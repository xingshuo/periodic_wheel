package periodic_wheel

import (
	"errors"
	"fmt"
	"time"
)

// 自己实现与框架绑定的优先队列
type Heapq struct {
	size int
	cap  int
	data []*tickFilter
}

func (self *Heapq) push(f *tickFilter) error {
	// 暂时不支持动态扩容
	if self.size >= self.cap {
		return errors.New("data size overflow")
	}
	self.size++
	self.data[self.size] = f
	self.shiftup(self.size)
	return nil
}

func (self *Heapq) top() *tickFilter {
	if self.size < 1 {
		return nil
	}
	return self.data[1]
}

func (self *Heapq) pop() *tickFilter {
	if self.size < 1 {
		return nil
	}
	top := self.data[1]
	self.data[1] = self.data[self.size]
	self.size--
	self.shiftdown(1)
	return top
}

func (self *Heapq) PopUntil(now time.Time) {
	nowTime := now.Unix()
	for self.size >= 1 {
		top := self.data[1]
		if top.expired {
			self.pop()
			continue
		}
		if nowTime < top.nextTickTime {
			break
		}
		top.onTick(nowTime)
		top.nextTickTime = top.getNextTickTime(now)
		self.shiftdown(1)
	}
}

func (self *Heapq) shiftup(k int) {
	v := self.data[k]
	for {
		c := k / 2
		if c <= 0 || self.data[c].nextTickTime <= v.nextTickTime {
			break
		}
		self.data[k] = self.data[c]
		k = c
	}
	self.data[k] = v
}

func (self *Heapq) shiftdown(k int) {
	v := self.data[k]
	for {
		c := k * 2
		if c > self.size {
			break
		}
		if c < self.size && self.data[c].nextTickTime > self.data[c+1].nextTickTime {
			c++
		}
		if v.nextTickTime <= self.data[c].nextTickTime {
			break
		}
		self.data[k] = self.data[c]
		k = c
	}
	self.data[k] = v
}

func (self *Heapq) debug() {
	fmt.Printf("========heapq debug========\n")
	for i := 1; i <= self.size; i++ {
		fmt.Printf("%dst %p %d\n", i, self.data[i], self.data[i].nextTickTime)
	}
}

func NewHeapq(cap int) *Heapq {
	q := &Heapq{
		size: 0,
		cap:  cap,
		data: make([]*tickFilter, cap+1),
	}
	return q
}
