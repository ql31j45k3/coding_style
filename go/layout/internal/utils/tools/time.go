package tools

import (
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

func GetStartTimeAndEndTime(startTimeStr, endTimeStr, timezone, timeFormat string) (time.Time, time.Time, error) {
	loc, err := time.LoadLocation(timezone)
	if err != nil {
		return time.Time{}, time.Time{}, fmt.Errorf("time.LoadLocation - %w", err)
	}

	startTime, err := time.ParseInLocation(timeFormat, startTimeStr, loc)
	if err != nil {
		return time.Time{}, time.Time{}, fmt.Errorf("time.ParseInLocation(startTimeStr) - %w", err)
	}

	endTime, err := time.ParseInLocation(timeFormat, endTimeStr, loc)
	if err != nil {
		return time.Time{}, time.Time{}, fmt.Errorf("time.ParseInLocation(endTimeStr) - %w", err)
	}

	return startTime, endTime, nil
}

func GetNowTime(timezone string) (time.Time, error) {
	t := time.Now()
	loadLocation, err := time.LoadLocation(timezone)
	if err != nil {
		return time.Time{}, fmt.Errorf("time.LoadLocation - %w", err)
	}

	return t.In(loadLocation), nil
}

func GetNowTimeStrAndFormat(timezone, layout string) (string, error) {
	nowTime, err := GetNowTime(timezone)
	if err != nil {
		return "", fmt.Errorf("GetNowTime - %w", err)
	}

	return nowTime.Format(layout), nil
}

func GetNowTimeAndFormat(timezone, layout string) (time.Time, error) {
	now, err := GetNowTime(timezone)
	if err != nil {
		return time.Time{}, fmt.Errorf("GetNowTime - %w", err)
	}

	loadLocation, err := time.LoadLocation(timezone)
	if err != nil {
		return time.Time{}, fmt.Errorf("time.LoadLocation - %w", err)
	}

	nowTime, err := time.ParseInLocation(layout, now.Format(layout), loadLocation)
	if err != nil {
		return time.Time{}, fmt.Errorf("time.ParseInLocation - %w", err)
	}

	return nowTime, nil
}

func GetNowTimestamp(timezone string) (int64, error) {
	nowTime, err := GetNowTime(timezone)
	if err != nil {
		return 0, fmt.Errorf("time.GetNowTimestamp - %w", err)
	}

	return GetTimeToTimestamp(nowTime), nil
}

func GetTimeStrToTimestamp(timeStr string, timezone string) (int64, error) {
	loc, err := time.LoadLocation(timezone)
	if err != nil {
		return 0, fmt.Errorf("time.LoadLocation - %w", err)
	}

	t, err := time.ParseInLocation(TimeFormatSecond, timeStr, loc)
	if err != nil {
		return 0, fmt.Errorf("time.ParseInLocation - %w", err)
	}

	return t.UnixNano() / 1e6, nil
}

func GetTimeStrToTimestampAndFormat(timeStr string, timeFormat, timezone string) (int64, error) {
	loc, err := time.LoadLocation(timezone)
	if err != nil {
		return 0, fmt.Errorf("time.LoadLocation - %w", err)
	}

	t, err := time.ParseInLocation(timeFormat, timeStr, loc)
	if err != nil {
		return 0, fmt.Errorf("time.ParseInLocation - %w", err)
	}

	return t.UnixNano() / 1e6, nil
}

func GetTimeStrToTime(timeStr string, timeFormat, timezone string) (time.Time, error) {
	loc, err := time.LoadLocation(timezone)
	if err != nil {
		return time.Time{}, fmt.Errorf("time.LoadLocation - %w", err)
	}

	t, err := time.ParseInLocation(timeFormat, timeStr, loc)
	if err != nil {
		return time.Time{}, fmt.Errorf("time.ParseInLocation - %w", err)
	}

	return t, nil
}

func GetTimeToTimeFormatAndTimezone(baseTime time.Time, timeFormat, timezone string) (time.Time, error) {
	loc, err := time.LoadLocation(timezone)
	if err != nil {
		return time.Time{}, fmt.Errorf("time.LoadLocation - %w", err)
	}

	baseTimeStr := baseTime.Format(timeFormat)
	t, err := time.ParseInLocation(timeFormat, baseTimeStr, loc)
	if err != nil {
		return time.Time{}, fmt.Errorf("time.ParseInLocation - %w", err)
	}

	return t, nil
}

func GetTimestampToStrFormat(timestamp int64, timezone, timeFormat string) (string, error) {
	tempTime, err := GetTimestampToTime(timestamp, timezone)
	if err != nil {
		return "", fmt.Errorf("GetTimestampToTime - %w", err)
	}

	timeStr := tempTime.Format(timeFormat)

	return timeStr, nil
}

func GetTodayTimestamp() (int64, int64, error) {
	nowTimeBasic, err := GetNowTime(TimezoneTaipei)
	if err != nil {
		return 0, 0, fmt.Errorf("GetNowTime - %w", err)
	}

	nowDate := nowTimeBasic.Format(TimeFormatDay) + " 00:00:00"

	// 時間轉換
	local, err := time.LoadLocation(TimezoneTaipei)
	if err != nil {
		return 0, 0, fmt.Errorf("time.LoadLocation - %w", err)
	}

	nowTime, err := time.ParseInLocation(TimeFormatSecond, nowDate, local)
	if err != nil {
		return 0, 0, fmt.Errorf("time.ParseInLocation - %w", err)
	}

	startTime, endTime := GetTimeStartAndEnd(nowTime)

	return startTime, endTime, nil
}

func GetYesterdayTimestamp(timeStr string) (int64, int64, error) {
	// 時間轉換
	local, err := time.LoadLocation(TimezoneTaipei)
	if err != nil {
		return 0, 0, fmt.Errorf("time.LoadLocation - %w", err)
	}

	timeEnd, err := time.ParseInLocation(TimeFormatSecond, timeStr, local)
	if err != nil {
		return 0, 0, fmt.Errorf("time.ParseInLocation - %w", err)
	}

	// 1天前
	timeStart := timeEnd.AddDate(0, 0, -1)
	startTime := timeStart.UnixNano() / 1e6
	endTime := timeEnd.UnixNano() / 1e6

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
