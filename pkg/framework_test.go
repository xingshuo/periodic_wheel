package periodic_wheel

import (
	"log"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

type TestUsrObj struct {
	TotalDay     int
	TotalWeek    int
	TotalMonth   int
	lastTickTime int64
}

func (u *TestUsrObj) OnNewDay(nowTime int64) {
	u.TotalDay++
	u.lastTickTime = nowTime
}

func (u *TestUsrObj) OnNewWeek(nowTime int64) {
	u.TotalWeek++
	u.lastTickTime = nowTime
}

func (u *TestUsrObj) OnNewMonth(nowTime int64) {
	u.TotalMonth++
	u.lastTickTime = nowTime
}

func TestPeriodicWheel(t *testing.T) {
	tm := time.Now()
	zone, offset := tm.Zone()
	log.Printf("zone:%v offset:%v\n", zone, offset)
	uobj1 := &TestUsrObj{}
	uobj2 := &TestUsrObj{}
	pw := NewPeriodicWheel()
	pw.PushDayFilter("TestUsrObj1.NewDay", 5, uobj1.OnNewDay)
	pw.PushWeekFilter("TestUsrObj1.NewWeek", 2, 5, uobj1.OnNewWeek)
	pw.PushMonthFilter("TestUsrObj1.NewMonth", 2, 5, uobj1.OnNewMonth)
	pw.PushDayFilter("TestUsrObj2.NewDay", 5, uobj2.OnNewDay)
	pw.PushWeekFilter("TestUsrObj2.NewWeek", 2, 5, uobj2.OnNewWeek)
	pw.PushMonthFilter("TestUsrObj2.NewMonth", 2, 5, uobj2.OnNewMonth)
	// 测试更新
	fakeTime := time.Date(2020, 12, 3, 10, 0, 0, 0, time.Local)
	pw.Update(fakeTime)
	assert.Equal(t, 1, uobj1.TotalDay)
	assert.Equal(t, 1, uobj1.TotalWeek)
	assert.Equal(t, 1, uobj1.TotalMonth)
	// 测试未过天
	fakeTime = time.Date(2020, 12, 4, 4, 0, 0, 0, time.Local)
	pw.Update(fakeTime)
	assert.Equal(t, 1, uobj1.TotalDay)
	assert.Equal(t, 1, uobj1.TotalWeek)
	assert.Equal(t, 1, uobj1.TotalMonth)
	// 测试过天(2020.12.4 星期五)
	fakeTime = time.Date(2020, 12, 4, 5, 0, 0, 0, time.Local)
	pw.Update(fakeTime)
	assert.Equal(t, 2, uobj1.TotalDay)
	assert.Equal(t, 1, uobj1.TotalWeek)
	assert.Equal(t, 1, uobj1.TotalMonth)
	// 测试未过周(2020.12.8 星期二)
	fakeTime = time.Date(2020, 12, 8, 4, 0, 0, 0, time.Local)
	pw.Update(fakeTime)
	assert.Equal(t, 3, uobj1.TotalDay)
	assert.Equal(t, 1, uobj1.TotalWeek)
	assert.Equal(t, 1, uobj1.TotalMonth)
	// 测试过周
	fakeTime = time.Date(2020, 12, 8, 5, 0, 0, 0, time.Local)
	pw.Update(fakeTime)
	assert.Equal(t, 4, uobj1.TotalDay)
	assert.Equal(t, 2, uobj1.TotalWeek)
	assert.Equal(t, 1, uobj1.TotalMonth)
	// 测试未过月
	fakeTime = time.Date(2021, 1, 2, 4, 0, 0, 0, time.Local)
	pw.Update(fakeTime)
	assert.Equal(t, 5, uobj1.TotalDay)
	assert.Equal(t, 3, uobj1.TotalWeek)
	assert.Equal(t, 1, uobj1.TotalMonth)
	// 测试过月
	fakeTime = time.Date(2021, 1, 2, 5, 0, 0, 0, time.Local)
	pw.Update(fakeTime)
	assert.Equal(t, 6, uobj1.TotalDay)
	assert.Equal(t, 3, uobj1.TotalWeek)
	assert.Equal(t, 2, uobj1.TotalMonth)
	// 测试删除
	pw.RemoveFilter("TestUsrObj1.NewDay")
	fakeTime = time.Date(2021, 5, 2, 5, 0, 0, 0, time.Local)
	pw.Update(fakeTime)
	assert.Equal(t, 6, uobj1.TotalDay)
	assert.Equal(t, 4, uobj1.TotalWeek)
	assert.Equal(t, 3, uobj1.TotalMonth)
	// 测试批量删除
	pw.BatchRemoveFilters("TestUsrObj1")
	fakeTime = time.Date(2022, 2, 2, 5, 0, 0, 0, time.Local)
	pw.Update(fakeTime)
	assert.Equal(t, 6, uobj1.TotalDay)
	assert.Equal(t, 4, uobj1.TotalWeek)
	assert.Equal(t, 3, uobj1.TotalMonth)
	// 测试注册filter的不同对象间独立性
	assert.Equal(t, 8, uobj2.TotalDay)
	assert.Equal(t, 5, uobj2.TotalWeek)
	assert.Equal(t, 4, uobj2.TotalMonth)
}
