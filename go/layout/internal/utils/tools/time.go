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
	Timezone string
	Layout   string
}

func (nt *NowTime) ToTime() (time.Time, error) {
	if IsEmpty(nt.Timezone) {
		return time.Time{}, errors.New("Timezone is not empty")
	}

	if IsEmpty(nt.Layout) {
		return time.Time{}, errors.New("Layout is not empty")
	}

	t := time.Now()
	loadLocation, err := time.LoadLocation(nt.Timezone)
	if err != nil {
		return time.Time{}, fmt.Errorf("time.LoadLocation - %w", err)
	}
	t = t.In(loadLocation)

	nowTime, err := time.ParseInLocation(nt.Layout, t.Format(nt.Layout), loadLocation)
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

	return nowTime.Format(nt.Layout), nil
}

type TimeStrValue struct {
	TimeStr string

	Timezone string
	Layout   string
}

func (tsv *TimeStrValue) ToTime() (time.Time, error) {
	t, err := ParseInLocation(tsv.TimeStr, tsv.Timezone, tsv.Layout)
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
	t, err := ParseInLocation(tsv.TimeStr, tsv.Timezone, TimeFormatSecond)
	if err != nil {
		return "", fmt.Errorf("ParseInLocation - %w", err)
	}

	timeStr := t.Format(tsv.Layout)
	return timeStr, nil
}

type TimeValue struct {
	BaseTime time.Time

	Timezone string
	Layout   string
}

func (tv *TimeValue) ToTime() (time.Time, error) {
	baseTimeStr := tv.BaseTime.Format(tv.Layout)
	t, err := ParseInLocation(baseTimeStr, tv.Timezone, tv.Layout)
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

	timeStr := t.Format(tv.Layout)
	return timeStr, nil
}

type TimestampValue struct {
	Timestamp int64

	Timezone string
	Layout   string
}

func (tv *TimestampValue) ToTime() (time.Time, error) {
	timeStr, err := tv.ToStr()
	if err != nil {
		return time.Time{}, fmt.Errorf("tv.ToStr - %w", err)
	}

	t, err := ParseInLocation(timeStr, tv.Timezone, tv.Layout)
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
	tempTime, err := TimestampConvTime(tv.Timestamp, tv.Timezone)
	if err != nil {
		return "", fmt.Errorf("GetTimestampToTime - %w", err)
	}

	timeStr := tempTime.Format(tv.Layout)
	return timeStr, nil
}

func GetTodayTimestampDefault() (int64, int64, error) {
	t := NowTime{
		Timezone: TimezoneTaipei,
		Layout:   TimeFormatSecond,
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

func GetMonthStartTimeAndEndTime(timeStr, timezone, layout string, years, months int) (time.Time, time.Time, error) {
	if IsEmpty(timezone) {
		timezone = TimezoneTaipei
	}

	if IsEmpty(layout) {
		layout = TimeFormatSecond
	}

	loc, err := time.LoadLocation(timezone)
	if err != nil {
		return time.Time{}, time.Time{}, fmt.Errorf("time.LoadLocation - %w", err)
	}

	basicTime, err := time.ParseInLocation(layout, timeStr, loc)
	if err != nil {
		return time.Time{}, time.Time{}, fmt.Errorf("time.ParseInLocation - %w", err)
	}

	year, month, _ := basicTime.AddDate(years, months, 0).Date()

	startTime := time.Date(year, month, 1, 0, 0, 0, 0, loc)
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
