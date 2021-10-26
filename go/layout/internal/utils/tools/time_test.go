package tools

import (
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestTimestampToMS(t *testing.T) {
	type args struct {
		timestamp int64
	}
	tests := []struct {
		name string
		args args
		want int64
	}{
		{
			name: "",
			args: args{
				timestamp: 1584806400,
			},
			want: 1584806400000,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := TimestampToMS(tt.args.timestamp)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestGetTimestampToTime(t *testing.T) {
	local, err := time.LoadLocation(TimezoneTaipei)
	if err != nil {
		t.Log(err)
		return
	}

	type args struct {
		timestamp int64
		timezone  string
	}
	tests := []struct {
		name    string
		args    args
		want    time.Time
		wantErr bool
	}{
		{
			name: "",
			args: args{
				timestamp: 1577808000000,
				timezone:  TimezoneTaipei,
			},
			want:    time.Date(2020, 1, 1, 0, 0, 0, 0, local),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetTimestampToTime(tt.args.timestamp, tt.args.timezone)
			if (err != nil) != tt.wantErr {
				assert.NoError(t, err, "GetTimestampToTime error = %v", err)
				return
			}

			assert.Equal(t, tt.want.String(), got.String())
		})
	}
}

func TestGetNowTimeStrAndFormat(t *testing.T) {
	nowTime, err := GetNowTime(TimezoneTaipei)
	if err != nil {
		t.Error(err)
		return
	}

	nowTimeStr := nowTime.Format(TimeFormatSecond)[:13] + ":00:00"

	type args struct {
		timezone string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "",
			args: args{
				TimezoneTaipei,
			},
			want:    nowTimeStr,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetNowTimeStrAndFormat(tt.args.timezone, TimeFormatHour)
			if (err != nil) != tt.wantErr {
				assert.NoError(t, err, "GetNowTimeStrAndFormat() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			assert.Equal(t, tt.want, got)
		})
	}
}

func TestGetNowTimeAndFormat(t *testing.T) {
	nowTime, err := GetNowTimeAndFormat(TimezoneTaipei, TimeFormatSecond)
	if err != nil {
		t.Error(err)
		return
	}

	type args struct {
		timezone string
		layout   string
	}
	tests := []struct {
		name    string
		args    args
		want    time.Time
		wantErr bool
	}{
		{
			name: "",
			args: args{
				timezone: TimezoneTaipei,
				layout:   TimeFormatSecond,
			},
			want:    nowTime,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetNowTimeAndFormat(tt.args.timezone, tt.args.layout)
			if (err != nil) != tt.wantErr {
				assert.NoError(t, err, "GetNowTimeAndFormat() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			assert.Equal(t, tt.want, got)
		})
	}
}

func TestGetNowTimestamp(t *testing.T) {
	nowTime, err := GetNowTime(TimezoneTaipei)
	if err != nil {
		t.Error(err)
		return
	}

	nowTimestamp := GetTimeToTimestamp(nowTime)

	type args struct {
		timezone string
	}
	tests := []struct {
		name    string
		args    args
		want    int64
		wantErr bool
	}{
		{
			name: "",
			args: args{
				timezone: TimezoneTaipei,
			},
			want:    nowTimestamp,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetNowTimestamp(tt.args.timezone)
			if (err != nil) != tt.wantErr {
				assert.NoError(t, err, "GetNowTimestamp() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			// 注意: 只比較前 10碼 到秒數比較即可
			wantStr := strconv.Itoa(int(tt.want))
			gotStr := strconv.Itoa(int(got))

			assert.Equal(t, wantStr[0:10], gotStr[0:10])
		})
	}
}

func TestGetTimeStartAndEnd(t *testing.T) {
	local, err := time.LoadLocation(TimezoneTaipei)
	if err != nil {
		t.Log(err)
		return
	}

	type args struct {
		nowTime time.Time
	}
	var tests = []struct {
		name  string
		args  args
		want  int64
		want1 int64
	}{
		{
			name: "nowTime= 2021-09-30 09:30:00",
			args: args{
				nowTime: time.Date(2021, 9, 30, 9, 30, 0, 0, local),
			},
			want:  1632965400000,
			want1: 1633051800000,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := GetTimeStartAndEnd(tt.args.nowTime)

			assert.Equal(t, tt.want, got)
			assert.Equal(t, tt.want1, got1)
		})
	}
}

func TestGetYesterdayTimestamp(t *testing.T) {
	type args struct {
		timeStr string
	}
	tests := []struct {
		name    string
		args    args
		want    int64
		want1   int64
		wantErr bool
	}{
		{
			name: "",
			args: args{
				timeStr: "2020-03-22 00:00:00",
			},
			want:    1584720000000,
			want1:   1584806400000,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1, err := GetYesterdayTimestamp(tt.args.timeStr)
			if (err != nil) != tt.wantErr {
				assert.NoError(t, err, "GetYesterdayTimestamp error = %v", err)
				return
			}
			assert.Equal(t, tt.want, got)
			assert.Equal(t, tt.want1, got1)
		})
	}
}

func TestGetMonthStartTimeAndEndTime(t *testing.T) {
	local, err := time.LoadLocation(TimezoneTaipei)
	if err != nil {
		t.Log(err)
		return
	}

	type args struct {
		timeStr string
		layout  string
		years   int
		months  int
	}
	tests := []struct {
		name    string
		args    args
		want    time.Time
		want1   time.Time
		wantErr bool
	}{
		{
			name: "",
			args: args{
				timeStr: "2020-03-22 00:00:00",
				layout:  TimeFormatSecond,
				years:   0,
				months:  -2,
			},
			want:    time.Date(2020, 1, 1, 0, 0, 0, 0, local),
			want1:   time.Date(2020, 1, 31, 0, 0, 0, 0, local),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1, err := GetMonthStartTimeAndEndTime(tt.args.timeStr, tt.args.layout, tt.args.years, tt.args.months)
			if (err != nil) != tt.wantErr {
				assert.NoError(t, err, "GetMonthStartTimeAndEndTime error = %v", err)
				return
			}

			assert.Equal(t, tt.want, got)
			assert.Equal(t, tt.want1, got1)
		})
	}
}

func TestGetTimeToTimestamp(t *testing.T) {
	local, err := time.LoadLocation(TimezoneTaipei)
	if err != nil {
		t.Log(err)
		return
	}

	type args struct {
		t time.Time
	}
	tests := []struct {
		name string
		args args
		want int64
	}{
		{
			name: "",
			args: args{
				t: time.Date(2020, 1, 1, 0, 0, 0, 0, local),
			},
			want: 1577808000000,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := GetTimeToTimestamp(tt.args.t)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestGetTimeStrToTimestamp(t *testing.T) {
	type args struct {
		timeStr  string
		timezone string
	}
	tests := []struct {
		name    string
		args    args
		want    int64
		wantErr bool
	}{
		{
			name: "",
			args: args{
				timeStr:  "2020-03-22 00:00:00",
				timezone: TimezoneTaipei,
			},
			want:    1584806400000,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetTimeStrToTimestamp(tt.args.timeStr, tt.args.timezone)
			if (err != nil) != tt.wantErr {
				assert.NoError(t, err, "GetTimeStrToTimestamp error = %v", err)
				return
			}
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestGetTimeStrToTimestampAndFormat(t *testing.T) {
	type args struct {
		timeStr    string
		timeFormat string
		timezone   string
	}
	tests := []struct {
		name    string
		args    args
		want    int64
		wantErr bool
	}{
		{
			name: "timeStr= 2021-09-30 09:30:00",
			args: args{
				timeStr:    "2021-09-30 09:30:00",
				timeFormat: TimeFormatSecond,
				timezone:   TimezoneTaipei,
			},
			want:    1632965400000,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetTimeStrToTimestampAndFormat(tt.args.timeStr, tt.args.timeFormat, tt.args.timezone)
			if (err != nil) != tt.wantErr {
				assert.NoError(t, err, "GetTimeStrToTimestampAndFormat() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			assert.Equal(t, tt.want, got)
		})
	}
}

func TestGetTimeStrToTime(t *testing.T) {
	local, err := time.LoadLocation(TimezoneTaipei)
	if err != nil {
		t.Log(err)
		return
	}

	type args struct {
		timeStr    string
		timeFormat string
		timezone   string
	}
	tests := []struct {
		name    string
		args    args
		want    time.Time
		wantErr bool
	}{
		{
			name: "timeStr= 2021-09-30 09:30:00",
			args: args{
				timeStr:    "2021-09-30 09:30:00",
				timeFormat: TimeFormatSecond,
				timezone:   TimezoneTaipei,
			},
			want:    time.Date(2021, 9, 30, 9, 30, 0, 0, local),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetTimeStrToTime(tt.args.timeStr, tt.args.timeFormat, tt.args.timezone)
			if (err != nil) != tt.wantErr {
				assert.NoError(t, err, "GetTimeStrToTime() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			assert.Equal(t, tt.want, got)
		})
	}
}

func TestGetTimeToTimeFormatAndTimezone(t *testing.T) {
	local, err := time.LoadLocation(TimezoneTaipei)
	if err != nil {
		t.Log(err)
		return
	}

	type args struct {
		baseTime   time.Time
		timeFormat string
		timezone   string
	}
	tests := []struct {
		name    string
		args    args
		want    time.Time
		wantErr bool
	}{
		{
			name: "",
			args: args{
				baseTime:   time.Date(2020, 1, 1, 0, 0, 0, 100, local),
				timeFormat: TimeFormatSecond,
				timezone:   TimezoneTaipei,
			},
			want:    time.Date(2020, 1, 1, 0, 0, 0, 0, local),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetTimeToTimeFormatAndTimezone(tt.args.baseTime, tt.args.timeFormat, tt.args.timezone)
			if (err != nil) != tt.wantErr {
				assert.NoError(t, err, "GetTimeToTimeFormatAndTimezone error = %v", err)
				return
			}
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestGetTimestampToStrFormat(t *testing.T) {
	type args struct {
		timestamp  int64
		timezone   string
		timeFormat string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "",
			args: args{
				timestamp:  1630861958000,
				timezone:   TimezoneTaipei,
				timeFormat: TimeFormatSecond,
			},
			want:    "2021-09-06 01:12:38",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetTimestampToStrFormat(tt.args.timestamp, tt.args.timezone, tt.args.timeFormat)
			if (err != nil) != tt.wantErr {
				assert.NoError(t, err, "GetTimestampToStrFormat() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			assert.Equal(t, tt.want, got)
		})
	}
}

func TestParseInLocation(t *testing.T) {
	local, err := time.LoadLocation(TimezoneTaipei)
	if err != nil {
		t.Log(err)
		return
	}

	type args struct {
		timeStr    string
		timezone   string
		timeFormat string
	}
	tests := []struct {
		name    string
		args    args
		want    time.Time
		wantErr bool
	}{
		{
			name: "",
			args: args{
				timeStr:    "2020-03-22 00:00:00",
				timezone:   TimezoneTaipei,
				timeFormat: TimeFormatSecond,
			},
			want:    time.Date(2020, 3, 22, 0, 0, 0, 0, local),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseInLocation(tt.args.timeStr, tt.args.timezone, tt.args.timeFormat)
			if (err != nil) != tt.wantErr {
				assert.NoError(t, err, "ParseInLocation error = %v", err)
				return
			}
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestGetStartTimeAndEndTime(t *testing.T) {
	local, err := time.LoadLocation(TimezoneTaipei)
	if err != nil {
		t.Log(err)
		return
	}

	type args struct {
		startTimeStr string
		endTimeStr   string
		timezone     string
		timeFormat   string
	}
	tests := []struct {
		name    string
		args    args
		want    time.Time
		want1   time.Time
		wantErr bool
	}{
		{
			name: "startTimeStr= 2021-09-30 09:30:00, endTimeStr= 2021-09-30 18:30:00, timezone= TimezoneTaipei, timeFormat= TimeFormatSecond",
			args: args{
				startTimeStr: "2021-09-30 09:30:00",
				endTimeStr:   "2021-09-30 18:30:00",
				timezone:     TimezoneTaipei,
				timeFormat:   TimeFormatSecond,
			},
			want:    time.Date(2021, 9, 30, 9, 30, 0, 0, local),
			want1:   time.Date(2021, 9, 30, 18, 30, 0, 0, local),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1, err := GetStartTimeAndEndTime(tt.args.startTimeStr, tt.args.endTimeStr, tt.args.timezone, tt.args.timeFormat)
			if (err != nil) != tt.wantErr {
				assert.NoError(t, err, "GetStartTimeAndEndTime() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			assert.Equal(t, tt.want, got)
			assert.Equal(t, tt.want1, got1)
		})
	}
}

func TestGetTimeSubDay(t *testing.T) {
	local, err := time.LoadLocation(TimezoneTaipei)
	if err != nil {
		t.Log(err)
		return
	}

	type args struct {
		t1 time.Time
		t2 time.Time
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{
			name: "one second but different day should return 1",
			args: args{
				t1: time.Date(2007, 1, 2, 23, 59, 59, 0, local),
				t2: time.Date(2007, 1, 3, 0, 0, 0, 0, local),
			},
			want: 1,
		},
		{
			name: "just one day should return 1, 2007",
			args: args{
				t1: time.Date(2007, 1, 2, 23, 59, 59, 0, local),
				t2: time.Date(2007, 1, 3, 23, 59, 59, 0, local),
			},
			want: 1,
		},
		{
			name: "just one day should return 1, 2017",
			args: args{
				t1: time.Date(2017, 9, 1, 10, 0, 0, 0, local),
				t2: time.Date(2017, 9, 2, 11, 0, 0, 0, local),
			},
			want: 1,
		},
		{
			name: "just one day should return 2",
			args: args{
				t1: time.Date(2007, 1, 2, 23, 59, 59, 0, local),
				t2: time.Date(2007, 1, 4, 0, 0, 0, 0, local),
			},
			want: 2,
		},
		{
			name: "just one day should return 3",
			args: args{
				t1: time.Date(2007, 1, 2, 0, 0, 0, 0, local),
				t2: time.Date(2007, 1, 5, 0, 0, 0, 0, local),
			},
			want: 3,
		},
		{
			name: "just one month:31 days should return 31",
			args: args{
				t1: time.Date(2007, 1, 2, 0, 0, 0, 0, local),
				t2: time.Date(2007, 2, 2, 0, 0, 0, 0, local),
			},
			want: 31,
		},
		{
			name: "just one month:29 days should return 29",
			args: args{
				t1: time.Date(2000, 2, 1, 0, 0, 0, 0, local),
				t2: time.Date(2000, 3, 1, 0, 0, 0, 0, local),
			},
			want: 29,
		},
		{
			name: "just one day: should return 1, Local",
			args: args{
				t1: time.Date(2018, 1, 9, 23, 59, 22, 100, time.Local),
				t2: time.Date(2018, 1, 10, 0, 0, 1, 100, time.Local),
			},
			want: 1,
		},
		{
			name: "just one day: should return 1, UTC",
			args: args{
				t1: time.Date(2018, 1, 9, 23, 59, 22, 100, local),
				t2: time.Date(2018, 1, 10, 0, 0, 1, 100, local),
			},
			want: 1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := GetTimeSubDay(tt.args.t2, tt.args.t1)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestIsTimeGte(t *testing.T) {
	local, err := time.LoadLocation(TimezoneTaipei)
	if err != nil {
		t.Log(err)
		return
	}

	type args struct {
		startTime time.Time
		t1        time.Time
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "",
			args: args{
				startTime: time.Date(2020, 1, 1, 0, 0, 0, 0, local),
				t1:        time.Date(2020, 1, 1, 0, 0, 0, 0, local),
			},
			want: true,
		},
		{
			name: "",
			args: args{
				startTime: time.Date(2020, 1, 1, 0, 0, 0, 0, local),
				t1:        time.Date(2020, 1, 1, 0, 0, 1, 0, local),
			},
			want: true,
		},
		{
			name: "",
			args: args{
				startTime: time.Date(2020, 1, 2, 0, 0, 0, 0, local),
				t1:        time.Date(2020, 1, 1, 23, 59, 59, 0, local),
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := IsTimeGte(tt.args.startTime, tt.args.t1)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestIsTimeLte(t *testing.T) {
	local, err := time.LoadLocation(TimezoneTaipei)
	if err != nil {
		t.Log(err)
		return
	}

	type args struct {
		t1 time.Time
		t2 time.Time
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "",
			args: args{
				t1: time.Date(2020, 1, 1, 0, 0, 0, 0, local),
				t2: time.Date(2020, 1, 1, 0, 0, 0, 0, local),
			},
			want: true,
		},
		{
			name: "",
			args: args{
				t1: time.Date(2020, 1, 1, 0, 0, 0, 0, local),
				t2: time.Date(2020, 1, 1, 0, 0, 1, 0, local),
			},
			want: false,
		},
		{
			name: "",
			args: args{
				t1: time.Date(2020, 1, 2, 0, 0, 0, 0, local),
				t2: time.Date(2020, 1, 1, 23, 59, 59, 0, local),
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := IsTimeLte(tt.args.t1, tt.args.t2)

			assert.Equal(t, tt.want, got)
		})
	}
}

func TestBetweenGteAndLt(t *testing.T) {
	local, err := time.LoadLocation(TimezoneTaipei)
	if err != nil {
		t.Log(err)
		return
	}

	type args struct {
		startTime time.Time
		endTime   time.Time
		t1        time.Time
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "2018/1/9 00:00:00 ~ 2018/1/9 23:59:22, 2018/1/9 23:59:22",
			args: args{
				startTime: time.Date(2018, 1, 9, 0, 0, 0, 100, local),
				endTime:   time.Date(2018, 1, 9, 23, 59, 22, 100, local),
				t1:        time.Date(2018, 1, 9, 23, 59, 22, 100, local),
			},
			want: false,
		},
		{
			name: "2018/1/9 00:00:00 ~ 2018/1/9 23:59:22, 2018/1/9 23:59:21",
			args: args{
				startTime: time.Date(2018, 1, 9, 0, 0, 0, 100, local),
				endTime:   time.Date(2018, 1, 9, 23, 59, 22, 100, local),
				t1:        time.Date(2018, 1, 9, 23, 59, 21, 100, local),
			},
			want: true,
		},
		{
			name: "2018/1/9 00:00:00 ~ 2018/1/9 23:59:22, 018/1/9 00:00:00",
			args: args{
				startTime: time.Date(2018, 1, 9, 0, 0, 0, 100, local),
				endTime:   time.Date(2018, 1, 9, 23, 59, 22, 100, local),
				t1:        time.Date(2018, 1, 9, 0, 0, 0, 100, local),
			},
			want: true,
		},
		{
			name: "2018/1/9 00:00:00 ~ 2018/1/9 23:59:22, 2018/1/8 23:59:22",
			args: args{
				startTime: time.Date(2018, 1, 9, 0, 0, 0, 100, local),
				endTime:   time.Date(2018, 1, 9, 23, 59, 22, 100, local),
				t1:        time.Date(2018, 1, 8, 23, 59, 22, 100, local),
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := BetweenGteAndLt(tt.args.startTime, tt.args.endTime, tt.args.t1)
			assert.Equal(t, tt.want, got)
		})
	}
}
