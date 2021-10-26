package tools

import (
	"errors"
	"fmt"
	"time"
)

func TimestampToMS(timestamp int64) int64 {
	return timestamp * 1000
}

// TimestampConvTime timestamp 帶入到毫秒 13碼
func TimestampConvTime(timestamp int64, timezone string) (time.Time, error) {
	loc, err := time.LoadLocation(timezone)
	if err != nil {
		return time.Time{}, fmt.Errorf("time.LoadLocation - %w", err)
	}

	return time.Unix(timestamp/1000, 0).In(loc), nil
}

func TimeConvTimestamp(t time.Time) int64 {
	return t.UnixNano() / 1e6
}

func ParseInLocation(timeStr, timezone, layout string) (time.Time, error) {
	loc, err := time.LoadLocation(timezone)
	if err != nil {
		return time.Time{}, fmt.Errorf("time.LoadLocation - %w", err)
	}

	t, err := time.ParseInLocation(layout, timeStr, loc)
	if err != nil {
		return time.Time{}, fmt.Errorf("time.ParseInLocation - %w", err)
	}

	return t, nil
}

func GetTimeStartAndEnd(nowTime time.Time) (int64, int64) {
	startTime := TimeConvTimestamp(nowTime)

	tempEndTime := nowTime.AddDate(0, 0, 1)
	endTime := TimeConvTimestamp(tempEndTime)

	return startTime, endTime
}

func GetStartTimeAndEndTime(startTimeStr, endTimeStr, timezone, layout string) (time.Time, time.Time, error) {
	loc, err := time.LoadLocation(timezone)
	if err != nil {
		return time.Time{}, time.Time{}, fmt.Errorf("time.LoadLocation - %w", err)
	}

	startTime, err := time.ParseInLocation(layout, startTimeStr, loc)
	if err != nil {
		return time.Time{}, time.Time{}, fmt.Errorf("time.ParseInLocation(startTimeStr) - %w", err)
	}

	endTime, err := time.ParseInLocation(layout, endTimeStr, loc)
	if err != nil {
		return time.Time{}, time.Time{}, fmt.Errorf("time.ParseInLocation(endTimeStr) - %w", err)
	}

	return startTime, endTime, nil
}

type NowTime struct {
	timezone string
	layout   string
}

func (nt *NowTime) ToTime() (time.Time, error) {
	if IsEmpty(nt.timezone) {
		return time.Time{}, errors.New("timezone is not empty")
	}

	if IsEmpty(nt.layout) {
		return time.Time{}, errors.New("layout is not empty")
	}

	t := time.Now()
	loadLocation, err := time.LoadLocation(nt.timezone)
	if err != nil {
		return time.Time{}, fmt.Errorf("time.LoadLocation - %w", err)
	}
	t = t.In(loadLocation)

	nowTime, err := time.ParseInLocation(nt.layout, t.Format(nt.layout), loadLocation)
	if err != nil {
		return time.Time{}, fmt.Errorf("time.ParseInLocation - %w", err)
	}

	return nowTime, nil
}

func (nt *NowTime) ToTimestamp() (int64, error) {
	nowTime, err := nt.ToTime()
	if err != nil {
		return 0, fmt.Errorf("nt.ToTime - %w", err)
	}

	return TimeConvTimestamp(nowTime), nil
}

func (nt *NowTime) ToStr() (string, error) {
	nowTime, err := nt.ToTime()
	if err != nil {
		return "", fmt.Errorf("nt.ToTime - %w", err)
	}

	return nowTime.Format(nt.layout), nil
}

type TimeStrValue struct {
	timeStr string

	timezone string
	layout   string
}

func (tsv *TimeStrValue) ToTime() (time.Time, error) {
	t, err := ParseInLocation(tsv.timeStr, tsv.timezone, tsv.layout)
	if err != nil {
		return time.Time{}, fmt.Errorf("ParseInLocation - %w", err)
	}

	return t, nil
}

func (tsv *TimeStrValue) ToTimestamp() (int64, error) {
	t, err := tsv.ToTime()
	if err != nil {
		return 0, fmt.Errorf("tsv.ToTime - %w", err)
	}

	return TimeConvTimestamp(t), nil
}

func (tsv *TimeStrValue) ToStr() (string, error) {
	t, err := ParseInLocation(tsv.timeStr, tsv.timezone, TimeFormatSecond)
	if err != nil {
		return "", fmt.Errorf("ParseInLocation - %w", err)
	}

	timeStr := t.Format(tsv.layout)
	return timeStr, nil
}

type TimeValue struct {
	baseTime time.Time

	timezone string
	layout   string
}

func (tv *TimeValue) ToTime() (time.Time, error) {
	baseTimeStr := tv.baseTime.Format(tv.layout)
	t, err := ParseInLocation(baseTimeStr, tv.timezone, tv.layout)
	if err != nil {
		return time.Time{}, fmt.Errorf("ParseInLocation - %w", err)
	}

	return t, nil
}

func (tv *TimeValue) ToTimestamp() (int64, error) {
	t, err := tv.ToTime()
	if err != nil {
		return 0, fmt.Errorf("tv.ToTime - %w", err)
	}

	return TimeConvTimestamp(t), nil
}

func (tv *TimeValue) ToStr() (string, error) {
	t, err := tv.ToTime()
	if err != nil {
		return "", fmt.Errorf("tv.ToTime - %w", err)
	}

	timeStr := t.Format(tv.layout)
	return timeStr, nil
}

type TimestampValue struct {
	timestamp int64

	timezone string
	layout   string
}

func (tv *TimestampValue) ToTime() (time.Time, error) {
	timeStr, err := tv.ToStr()
	if err != nil {
		return time.Time{}, fmt.Errorf("tv.ToStr - %w", err)
	}

	t, err := ParseInLocation(timeStr, tv.timezone, tv.layout)
	if err != nil {
		return time.Time{}, fmt.Errorf("ParseInLocation - %w", err)
	}

	return t, nil
}

func (tv *TimestampValue) ToTimestamp() (int64, error) {
	t, err := tv.ToTime()
	if err != nil {
		return 0, fmt.Errorf("tv.ToTime - %w", err)
	}

	return TimeConvTimestamp(t), nil
}

func (tv *TimestampValue) ToStr() (string, error) {
	tempTime, err := TimestampConvTime(tv.timestamp, tv.timezone)
	if err != nil {
		return "", fmt.Errorf("GetTimestampToTime - %w", err)
	}

	timeStr := tempTime.Format(tv.layout)
	return timeStr, nil
}

func GetTodayTimestampDefault() (int64, int64, error) {
	t := NowTime{
		timezone: TimezoneTaipei,
		layout:   TimeFormatSecond,
	}

	nowTimeBasic, err := t.ToTime()
	if err != nil {
		return 0, 0, fmt.Errorf("t.ToTime() - %w", err)
	}

	return GetTodayTimestamp(nowTimeBasic, "", "")
}

func GetTodayTimestamp(nowTimeBasic time.Time, timezone, layout string) (int64, int64, error) {
	if IsEmpty(timezone) {
		timezone = TimezoneTaipei
	}

	if IsEmpty(layout) {
		layout = TimeFormatSecond
	}

	nowDate := nowTimeBasic.Format(TimeFormatDay) + " 00:00:00"

	nowTime, err := ParseInLocation(nowDate, timezone, layout)
	if err != nil {
		return 0, 0, fmt.Errorf("ParseInLocation - %w", err)
	}

	startTime, endTime := GetTimeStartAndEnd(nowTime)

	return startTime, endTime, nil
}

func GetYesterdayTimestampDefault(timeStr string) (int64, int64, error) {
	return GetYesterdayTimestamp(timeStr, "", "")
}

func GetYesterdayTimestamp(timeStr, timezone, layout string) (int64, int64, error) {
	if IsEmpty(timezone) {
		timezone = TimezoneTaipei
	}

	if IsEmpty(layout) {
		layout = TimeFormatSecond
	}

	timeEnd, err := ParseInLocation(timeStr, timezone, layout)
	if err != nil {
		return 0, 0, fmt.Errorf("ParseInLocation - %w", err)
	}

	// 1天前
	timeStart := timeEnd.AddDate(0, 0, -1)
	startTime := TimeConvTimestamp(timeStart)
	endTime := TimeConvTimestamp(timeEnd)

	return startTime, endTime, nil
}

func GetMonthStartTimeAndEndTime(timeStr, layout string, years, months int) (time.Time, time.Time, error) {
	// 時間轉換
	local, err := time.LoadLocation(TimezoneTaipei)
	if err != nil {
		return time.Time{}, time.Time{}, fmt.Errorf("time.LoadLocation - %w", err)
	}

	basicTime, err := time.ParseInLocation(layout, timeStr, local)
	if err != nil {
		return time.Time{}, time.Time{}, fmt.Errorf("time.ParseInLocation - %w", err)
	}

	year, month, _ := basicTime.AddDate(years, months, 0).Date()

	startTime := time.Date(year, month, 1, 0, 0, 0, 0, local)
	endTime := startTime.AddDate(0, 1, -1)

	return startTime, endTime, nil
}

func GetTimeSubDay(t1, t2 time.Time) int {
	t1 = time.Date(t1.Year(), t1.Month(), t1.Day(), 0, 0, 0, 0, t1.Location())
	t2 = time.Date(t2.Year(), t2.Month(), t2.Day(), 0, 0, 0, 0, t2.Location())

	return int(t1.Sub(t2).Hours() / 24)
}

func IsTimeGte(t1, t2 time.Time) bool {
	if t1.Equal(t2) {
		return true
	}

	if t1.Before(t2) {
		return true
	}

	return false
}

func IsTimeLte(t1, t2 time.Time) bool {
	if t1.Equal(t2) {
		return true
	}

	if t1.After(t2) {
		return true
	}

	return false
}

func BetweenGteAndLt(startTime, endTime, t1 time.Time) bool {
	if startTime.Equal(t1) {
		return true
	}

	if startTime.Before(t1) && endTime.After(t1) {
		return true
	}

	return false
}
