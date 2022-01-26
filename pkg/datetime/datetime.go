package datetime

import (
	"time"
)

// 日期转时间戳
func DateToTime(date string) int64 {
	timeFormat := "2006-01-02 15:04:05"

	// 时区
	Loc, _ := time.LoadLocation("Local")

	stamp, _ := time.ParseInLocation(timeFormat, date, Loc)

	return stamp.Unix()
}
