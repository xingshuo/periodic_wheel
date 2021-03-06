package periodic_wheel

import (
	"time"
)

// 每天N点 0 ~ 23
func GetNextDayTickTime(now time.Time, hour int) int64 {
	if now.Hour() < hour {
		return time.Date(now.Year(), now.Month(), now.Day(), hour, 0, 0, 0, now.Location()).Unix()
	}
	return time.Date(now.Year(), now.Month(), now.Day()+1, hour, 0, 0, 0, now.Location()).Unix()
}

// 每周X N点
// weekday: 1 ~ 7
// hour: 0 ~ 23
func GetNextWeekTickTime(now time.Time, weekday, hour int) int64 {
	curWeekDay := int(now.Weekday())
	curHour := now.Hour()
	if curWeekDay == 0 {
		curWeekDay = 7
	}
	if curWeekDay < weekday || (curWeekDay == weekday && curHour < hour) {
		dt := (weekday-curWeekDay)*24*3600 + (hour-curHour)*3600
		return time.Date(now.Year(), now.Month(), now.Day(), curHour, 0, 0, 0, now.Location()).Unix() + int64(dt)
	}
	dt := (weekday-curWeekDay)*24*3600 + (hour-curHour)*3600 + 7*24*3600
	return time.Date(now.Year(), now.Month(), now.Day(), curHour, 0, 0, 0, now.Location()).Unix() + int64(dt)
}

// 每月X号 N点
// day: 1 ~ 28 or -1 ~ -28
// hour: 0 ~ 23
func GetNextMonthTickTime(now time.Time, day, hour int) int64 {
	if day == 0 {
		day = -1
	}
	curDay := now.Day()
	curHour := now.Hour()
	if day > 0 {
		if curDay < day || (curDay == day && curHour < hour) {
			dt := (day-curDay)*24*3600 + (hour-curHour)*3600
			return time.Date(now.Year(), now.Month(), curDay, curHour, 0, 0, 0, now.Location()).Unix() + int64(dt)
		}
		return time.Date(now.Year(), now.Month(), day, hour, 0, 0, 0, now.Location()).AddDate(0, 1, 0).Unix()
	}
	curMonthLastDate := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location()).AddDate(0, 1, -1)
	nextDay := curMonthLastDate.Day() + day + 1
	if curDay < nextDay || (curDay == nextDay && curHour < hour) {
		dt := (nextDay-curDay)*24*3600 + (hour-curHour)*3600
		return time.Date(now.Year(), now.Month(), curDay, curHour, 0, 0, 0, now.Location()).Unix() + int64(dt)
	}
	return time.Date(now.Year(), now.Month(), 1, hour, 0, 0, 0, now.Location()).AddDate(0, 2, day).Unix()
}
