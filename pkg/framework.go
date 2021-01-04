package periodic_wheel

import (
	"fmt"
	"strings"
	"time"
)

const (
	MAX_FILTER_NUM = 1024
)

type tickFilter struct {
	// 获取下一次触发时间
	getNextTickTime func(now time.Time) int64
	// 业务tick逻辑
	onTick func(nowTime int64)
	// 下一次触发时间
	nextTickTime int64
	// 是否失效
	expired bool
}

type PeriodicWheel struct {
	queue   *Heapq
	filters map[string]*tickFilter
}

func (self *PeriodicWheel) Update(now time.Time) {
	self.queue.PopUntil(now)
}

// hour: 0 ~ 23
func (self *PeriodicWheel) PushDayFilter(key string, hour int, onTickFn func(int64)) error {
	if hour < 0 || hour > 23 {
		return fmt.Errorf("hour param fmt err")
	}
	if key == "" {
		return fmt.Errorf("invalid key")
	}
	oldf := self.filters[key]
	if oldf != nil {
		oldf.expired = true
	}
	getNextTmFn := func(now time.Time) int64 {
		return GetNextDayTickTime(now, hour)
	}
	f := &tickFilter{
		getNextTickTime: getNextTmFn,
		onTick:          onTickFn,
		nextTickTime:    getNextTmFn(time.Unix(0, 0)),
	}
	self.filters[key] = f
	return self.queue.push(f)
}

// weekday: 1 ~ 7
// hour: 0 ~ 23
func (self *PeriodicWheel) PushWeekFilter(key string, weekday, hour int, onTickFn func(int64)) error {
	if weekday <= 0 || weekday > 7 {
		return fmt.Errorf("weekday param fmt err")
	}
	if hour < 0 || hour > 23 {
		return fmt.Errorf("hour param fmt err")
	}
	if key == "" {
		return fmt.Errorf("invalid key")
	}
	oldf := self.filters[key]
	if oldf != nil {
		oldf.expired = true
	}
	getNextTmFn := func(now time.Time) int64 {
		return GetNextWeekTickTime(now, weekday, hour)
	}
	f := &tickFilter{
		getNextTickTime: getNextTmFn,
		onTick:          onTickFn,
		nextTickTime:    getNextTmFn(time.Unix(0, 0)),
	}
	self.filters[key] = f
	return self.queue.push(f)
}

// day: 1 ~ 28 or -1 ~ -28
// hour: 0 ~ 23
func (self *PeriodicWheel) PushMonthFilter(key string, day, hour int, onTickFn func(int64)) error {
	if !(1 <= day && day <= 28) && !(-28 <= day && day <= -1) {
		return fmt.Errorf("day param fmt err")
	}
	if hour < 0 || hour > 23 {
		return fmt.Errorf("hour param fmt err")
	}
	if key == "" {
		return fmt.Errorf("invalid key")
	}
	oldf := self.filters[key]
	if oldf != nil {
		oldf.expired = true
	}
	getNextTmFn := func(now time.Time) int64 {
		return GetNextMonthTickTime(now, day, hour)
	}
	f := &tickFilter{
		getNextTickTime: getNextTmFn,
		onTick:          onTickFn,
		nextTickTime:    getNextTmFn(time.Unix(0, 0)),
	}
	self.filters[key] = f
	return self.queue.push(f)
}

func (self *PeriodicWheel) PushCustomizedFilter(key string, getNextTmFn func(time.Time) int64, onTickFn func(int64)) error {
	oldf := self.filters[key]
	if oldf != nil {
		oldf.expired = true
	}
	f := &tickFilter{
		getNextTickTime: getNextTmFn,
		onTick:          onTickFn,
		nextTickTime:    getNextTmFn(time.Unix(0, 0)),
	}
	self.filters[key] = f
	return self.queue.push(f)
}

func (self *PeriodicWheel) RemoveFilter(key string) {
	f := self.filters[key]
	if f != nil {
		f.expired = true
	}
	delete(self.filters, key)
}

func (self *PeriodicWheel) BatchRemoveFilters(keyPrefix string) {
	for key, f := range self.filters {
		if strings.HasPrefix(key, keyPrefix) {
			f.expired = true
			delete(self.filters, key)
		}
	}
}

func NewPeriodicWheel() *PeriodicWheel {
	pw := &PeriodicWheel{
		filters: make(map[string]*tickFilter),
		queue:   NewHeapq(MAX_FILTER_NUM),
	}
	return pw
}
