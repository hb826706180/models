package models

import (
	"fmt"
	"strings"
	"time"
)

func GetDate(date time.Time) string {
	return date.Format("2006-01-02 15:04:05")
}

// GetUnix 获取时间戳() (unix int)
// 参数 n: 10-10位,13-13位,19-19位 ,可空 默认13位
// 返回值: 时间戳
func GetUnix(n int) int64 {
	now := time.Now()
	//10位数的时间戳是以 秒 为单位；
	//13位数的时间戳是以 毫秒 为单位；
	//19位数的时间戳是以 纳秒 为单位；
	//fmt.Println("时间戳(秒)：", now.Unix())
	//fmt.Println("时间戳(纳秒):", now.UnixNano())
	//fmt.Println("时间戳(毫秒)：", now.UnixNano()/1e6)
	//fmt.Println("时间戳(纳秒转毫秒)：", now.UnixNano()/1e9)

	if n == 10 {
		return now.Unix()
	} else if n == 13 {
		return now.UnixNano() / 1e6
	} else if n == 19 {
		return now.UnixNano() / 1e9
	}
	return now.UnixNano() / 1e6
}

// 时间日期()(Date string)
func UnixDate() (Date string) {
	t := time.Unix(time.Now().Unix(), 0)
	Date = t.Format("2006-01-02 15:04:05")
	return
}

// 时间日期() (day string)
func UnixDay() (day string) {
	day = time.Now().Format("20060102")
	return
}

// 时间戳转日期(timestamp int) (ret string)
func UnixToTime(timestamp int) (ret string) {
	fmt.Println(timestamp)
	t := time.Unix(int64(timestamp), 0)
	ret = t.Format("2006-01-02 15:04:05")
	return
}

// StrTimeToUnit 日期转时间戳 (Mon, 02 Jan 2006 15:04:05 MST)
func StrTimeToUnit(props string) int64 {
	times, _ := time.Parse("Mon, 02 Jan 2006 15:04:05 MST", props)
	timestamp := times.Unix()
	return timestamp
}

type String struct{}

// 删除所有指定的字符(timestamp int) (ret string)
func (a String) DeleteAllByte(s string, old string) string {
	s = strings.ReplaceAll(s, old, "")
	return s
}
