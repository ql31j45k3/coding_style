package tools

import (
	"fmt"
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

func TestTimestampConvTime(t *testing.T) {
	loc, err := time.LoadLocation(TimezoneTaipei)
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
			want:    time.Date(2020, 1, 1, 0, 0, 0, 0, loc),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := TimestampConvTime(tt.args.timestamp, tt.args.timezone)
			if (err != nil) != tt.wantErr {
				assert.NoError(t, err, "TimestampConvTime error = %v", err)
				return
			}

			assert.Equal(t, tt.want.String(), got.String())
		})
	}
}

func TestTimeConvTimestamp(t *testing.T) {
	loc, err := time.LoadLocation(TimezoneTaipei)
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
				t: time.Date(2020, 1, 1, 0, 0, 0, 0, loc),
			},
			want: 1577808000000,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := TimeConvTimestamp(tt.args.t)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestParseInLocation(t *testing.T) {
	loc, err := time.LoadLocation(TimezoneTaipei)
	if err != nil {
		t.Log(err)
		return
	}

	type args struct {
		timeStr  string
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
				timeStr:  "2020-03-22 00:00:00",
				timezone: TimezoneTaipei,
				layout:   TimeFormatSecond,
			},
			want:    time.Date(2020, 3, 22, 0, 0, 0, 0, loc),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseInLocation(tt.args.timeStr, tt.args.timezone, tt.args.layout)
			if (err != nil) != tt.wantErr {
				assert.NoError(t, err, "ParseInLocation error = %v", err)
				return
			}
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestGetTimeStartAndEnd(t *testing.T) {
	loc, err := time.LoadLocation(TimezoneTaipei)
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
				nowTime: time.Date(2021, 9, 30, 9, 30, 0, 0, loc),
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

func TestGetStartTimeAndEndTime(t *testing.T) {
	loc, err := time.LoadLocation(TimezoneTaipei)
	if err != nil {
		t.Log(err)
		return
	}

	type args struct {
		startTimeStr string
		endTimeStr   string
		timezone     string
		layout       string
	}
	tests := []struct {
		name    string
		args    args
		want    time.Time
		want1   time.Time
		wantErr bool
	}{
		{
			name: "startTimeStr= 2021-09-30 09:30:00, endTimeStr= 2021-09-30 18:30:00, Timezone= TimezoneTaipei, Layout= TimeFormatSecond",
			args: args{
				startTimeStr: "2021-09-30 09:30:00",
				endTimeStr:   "2021-09-30 18:30:00",
				timezone:     TimezoneTaipei,
				layout:       TimeFormatSecond,
			},
			want:    time.Date(2021, 9, 30, 9, 30, 0, 0, loc),
			want1:   time.Date(2021, 9, 30, 18, 30, 0, 0, loc),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1, err := GetStartTimeAndEndTime(tt.args.startTimeStr, tt.args.endTimeStr, tt.args.timezone, tt.args.layout)
			if (err != nil) != tt.wantErr {
				assert.NoError(t, err, "GetStartTimeAndEndTime() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			assert.Equal(t, tt.want, got)
			assert.Equal(t, tt.want1, got1)
		})
	}
}

func getNowTime(timezone, layout string) (time.Time, error) {
	t := time.Now()
	loadLocation, err := time.LoadLocation(timezone)
	if err != nil {
		return time.Time{}, fmt.Errorf("time.LoadLocation - %w", err)
	}
	t = t.In(loadLocation)

	nowTime, err := time.ParseInLocation(layout, t.Format(layout), loadLocation)
	if err != nil {
		return time.Time{}, fmt.Errorf("time.ParseInLocation - %w", err)
	}

	return nowTime, nil
}

func TestNowTime_ToTime(t *testing.T) {
	nowTime, err := getNowTime(TimezoneTaipei, TimeFormatSecond)
	if err != nil {
		t.Log(err)
		return
	}

	type fields struct {
		timezone string
		layout   string
	}
	tests := []struct {
		name    string
		fields  fields
		want    time.Time
		wantErr bool
	}{
		{
			name: "nowTime Layout= TimeFormatSecond",
			fields: fields{
				timezone: TimezoneTaipei,
				layout:   TimeFormatSecond,
			},
			want:    nowTime,
			wantErr: false,
		},
		{
			name: "nowTime Layout= TimeFormatDay",
			fields: fields{
				timezone: TimezoneTaipei,
				layout:   TimeFormatDay,
			},
			want:    nowTime,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			nt := &NowTime{
				Timezone: tt.fields.timezone,
				Layout:   tt.fields.layout,
			}
			got, err := nt.ToTime()
			if (err != nil) != tt.wantErr {
				t.Errorf("ToTime() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			assert.Equal(t, tt.want.Format(tt.fields.layout), got.Format(tt.fields.layout))
		})
	}
}

func TestNowTime_ToTimestamp(t *testing.T) {
	nowTime, err := getNowTime(TimezoneTaipei, TimeFormatSecond)
	if err != nil {
		t.Log(err)
		return
	}

	nowTime2, err := getNowTime(TimezoneTaipei, TimeFormatDay)
	if err != nil {
		t.Log(err)
		return
	}

	type fields struct {
		timezone string
		layout   string
	}
	tests := []struct {
		name    string
		fields  fields
		want    int64
		wantErr bool
	}{
		{
			name: "nowTime Layout= TimeFormatSecond",
			fields: fields{
				timezone: TimezoneTaipei,
				layout:   TimeFormatSecond,
			},
			want:    TimeConvTimestamp(nowTime),
			wantErr: false,
		},
		{
			name: "nowTime Layout= TimeFormatDay",
			fields: fields{
				timezone: TimezoneTaipei,
				layout:   TimeFormatDay,
			},
			want:    TimeConvTimestamp(nowTime2),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			nt := &NowTime{
				Timezone: tt.fields.timezone,
				Layout:   tt.fields.layout,
			}
			got, err := nt.ToTimestamp()
			if (err != nil) != tt.wantErr {
				t.Errorf("ToTimestamp() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestNowTime_ToStr(t *testing.T) {
	nowTime, err := getNowTime(TimezoneTaipei, TimeFormatSecond)
	if err != nil {
		t.Log(err)
		return
	}

	type fields struct {
		timezone string
		layout   string
	}
	tests := []struct {
		name    string
		fields  fields
		want    string
		wantErr bool
	}{
		{
			name: "nowTime Layout= TimeFormatSecond",
			fields: fields{
				timezone: TimezoneTaipei,
				layout:   TimeFormatSecond,
			},
			want:    nowTime.Format(TimeFormatSecond),
			wantErr: false,
		},
		{
			name: "nowTime Layout= TimeFormatDay",
			fields: fields{
				timezone: TimezoneTaipei,
				layout:   TimeFormatDay,
			},
			want:    nowTime.Format(TimeFormatDay),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			nt := &NowTime{
				Timezone: tt.fields.timezone,
				Layout:   tt.fields.layout,
			}
			got, err := nt.ToStr()
			if (err != nil) != tt.wantErr {
				t.Errorf("ToStr() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestTimeStrValue_ToTime(t *testing.T) {
	loc, err := time.LoadLocation(TimezoneTaipei)
	if err != nil {
		t.Log(err)
		return
	}

	type fields struct {
		timeStr  string
		timezone string
		layout   string
	}
	tests := []struct {
		name    string
		fields  fields
		want    time.Time
		wantErr bool
	}{
		{
			name: "TimeStr = 2020-03-22 00:00:00",
			fields: fields{
				timeStr:  "2020-03-22 00:00:00",
				timezone: TimezoneTaipei,
				layout:   TimeFormatSecond,
			},
			want:    time.Date(2020, 3, 22, 0, 0, 0, 0, loc),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tsv := &TimeStrValue{
				TimeStr:  tt.fields.timeStr,
				Timezone: tt.fields.timezone,
				Layout:   tt.fields.layout,
			}
			got, err := tsv.ToTime()
			if (err != nil) != tt.wantErr {
				t.Errorf("ToTime() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestTimeStrValue_ToTimestamp(t *testing.T) {
	type fields struct {
		timeStr  string
		timezone string
		layout   string
	}
	tests := []struct {
		name    string
		fields  fields
		want    int64
		wantErr bool
	}{
		{
			name: "TimeStr = 2020-03-22 00:00:00",
			fields: fields{
				timeStr:  "2020-03-22 00:00:00",
				timezone: TimezoneTaipei,
				layout:   TimeFormatSecond,
			},
			want:    1584806400000,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tsv := &TimeStrValue{
				TimeStr:  tt.fields.timeStr,
				Timezone: tt.fields.timezone,
				Layout:   tt.fields.layout,
			}
			got, err := tsv.ToTimestamp()
			if (err != nil) != tt.wantErr {
				t.Errorf("ToTimestamp() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestTimeStrValue_ToStr(t *testing.T) {
	type fields struct {
		timeStr  string
		timezone string
		layout   string
	}
	tests := []struct {
		name    string
		fields  fields
		want    string
		wantErr bool
	}{
		{
			name: "TimeStr = 2020-03-22 01:02:03, Layout= TimeFormatDay",
			fields: fields{
				timeStr:  "2020-03-22 01:02:03",
				timezone: TimezoneTaipei,
				layout:   TimeFormatDay,
			},
			want:    "2020-03-22",
			wantErr: false,
		},
		{
			name: "TimeStr = 2020-03-22 01:02:03, Layout= TimeFormatHour",
			fields: fields{
				timeStr:  "2020-03-22 01:02:03",
				timezone: TimezoneTaipei,
				layout:   TimeFormatHour,
			},
			want:    "2020-03-22 01:00:00",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tsv := &TimeStrValue{
				TimeStr:  tt.fields.timeStr,
				Timezone: tt.fields.timezone,
				Layout:   tt.fields.layout,
			}
			got, err := tsv.ToStr()
			if (err != nil) != tt.wantErr {
				t.Errorf("ToStr() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestTimeValue_ToTime(t *testing.T) {
	loc, err := time.LoadLocation(TimezoneTaipei)
	if err != nil {
		t.Log(err)
		return
	}

	type fields struct {
		baseTime time.Time
		timezone string
		layout   string
	}
	tests := []struct {
		name    string
		fields  fields
		want    time.Time
		wantErr bool
	}{
		{
			name: "BaseTime = 2020-03-22 01:02:03, Layout= TimeFormatHour",
			fields: fields{
				baseTime: time.Date(2020, 3, 22, 1, 2, 3, 0, loc),
				timezone: TimezoneTaipei,
				layout:   TimeFormatHour,
			},
			want:    time.Date(2020, 3, 22, 1, 0, 0, 0, loc),
			wantErr: false,
		},
		{
			name: "BaseTime = 2020-03-22 01:02:03, Layout= TimeFormatMonth",
			fields: fields{
				baseTime: time.Date(2020, 3, 22, 1, 2, 3, 0, loc),
				timezone: TimezoneTaipei,
				layout:   TimeFormatMonth,
			},
			want:    time.Date(2020, 3, 1, 0, 0, 0, 0, loc),
			wantErr: false,
		},
		{
			name: "BaseTime = 2020-03-22 01:02:03, Layout= TimeFormatHour2",
			fields: fields{
				baseTime: time.Date(2020, 3, 22, 1, 2, 3, 0, loc),
				timezone: TimezoneTaipei,
				layout:   TimeFormatHour2,
			},
			want:    time.Date(2020, 3, 22, 1, 0, 0, 0, loc),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tv := &TimeValue{
				BaseTime: tt.fields.baseTime,
				Timezone: tt.fields.timezone,
				Layout:   tt.fields.layout,
			}
			got, err := tv.ToTime()
			if (err != nil) != tt.wantErr {
				t.Errorf("ToTime() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			fmt.Println(tt.want, got)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestTimeValue_ToTimestamp(t *testing.T) {
	loc, err := time.LoadLocation(TimezoneTaipei)
	if err != nil {
		t.Log(err)
		return
	}

	type fields struct {
		baseTime time.Time
		timezone string
		layout   string
	}
	tests := []struct {
		name    string
		fields  fields
		want    int64
		wantErr bool
	}{
		{
			name: "BaseTime = 2020-03-22 01:02:03, Layout= TimeFormatHour",
			fields: fields{
				baseTime: time.Date(2020, 3, 22, 1, 2, 3, 0, loc),
				timezone: TimezoneTaipei,
				layout:   TimeFormatHour,
			},
			want:    TimeConvTimestamp(time.Date(2020, 3, 22, 1, 0, 0, 0, loc)),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tv := &TimeValue{
				BaseTime: tt.fields.baseTime,
				Timezone: tt.fields.timezone,
				Layout:   tt.fields.layout,
			}
			got, err := tv.ToTimestamp()
			if (err != nil) != tt.wantErr {
				t.Errorf("ToTimestamp() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestTimeValue_ToStr(t *testing.T) {
	loc, err := time.LoadLocation(TimezoneTaipei)
	if err != nil {
		t.Log(err)
		return
	}

	type fields struct {
		baseTime time.Time
		timezone string
		layout   string
	}
	tests := []struct {
		name    string
		fields  fields
		want    string
		wantErr bool
	}{
		{
			name: "BaseTime = 2020-03-22 01:02:03, Layout= TimeFormatHour",
			fields: fields{
				baseTime: time.Date(2020, 3, 22, 1, 2, 3, 0, loc),
				timezone: TimezoneTaipei,
				layout:   TimeFormatHour,
			},
			want:    time.Date(2020, 3, 22, 1, 0, 0, 0, loc).Format(TimeFormatHour),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tv := &TimeValue{
				BaseTime: tt.fields.baseTime,
				Timezone: tt.fields.timezone,
				Layout:   tt.fields.layout,
			}
			got, err := tv.ToStr()
			if (err != nil) != tt.wantErr {
				t.Errorf("ToStr() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestTimestampValue_ToTime(t *testing.T) {
	loc, err := time.LoadLocation(TimezoneTaipei)
	if err != nil {
		t.Log(err)
		return
	}

	type fields struct {
		timestamp int64
		timezone  string
		layout    string
	}
	tests := []struct {
		name    string
		fields  fields
		want    time.Time
		wantErr bool
	}{
		{
			name: "BaseTime = 2020-03-22 01:02:03, Layout= TimeFormatHour",
			fields: fields{
				timestamp: 1584810123000,
				timezone:  TimezoneTaipei,
				layout:    TimeFormatHour,
			},
			want:    time.Date(2020, 3, 22, 1, 0, 0, 0, loc),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tv := &TimestampValue{
				Timestamp: tt.fields.timestamp,
				Timezone:  tt.fields.timezone,
				Layout:    tt.fields.layout,
			}
			got, err := tv.ToTime()
			if (err != nil) != tt.wantErr {
				t.Errorf("ToTime() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestTimestampValue_ToTimestamp(t *testing.T) {
	type fields struct {
		timestamp int64
		timezone  string
		layout    string
	}
	tests := []struct {
		name    string
		fields  fields
		want    int64
		wantErr bool
	}{
		{
			name: "BaseTime = 2020-03-22 01:02:03, Layout= TimeFormatHour",
			fields: fields{
				timestamp: 1584810123000,
				timezone:  TimezoneTaipei,
				layout:    TimeFormatHour,
			},
			want:    1584810000000,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tv := &TimestampValue{
				Timestamp: tt.fields.timestamp,
				Timezone:  tt.fields.timezone,
				Layout:    tt.fields.layout,
			}
			got, err := tv.ToTimestamp()
			if (err != nil) != tt.wantErr {
				t.Errorf("ToTimestamp() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestTimestampValue_ToStr(t *testing.T) {
	type fields struct {
		timestamp int64
		timezone  string
		layout    string
	}
	tests := []struct {
		name    string
		fields  fields
		want    string
		wantErr bool
	}{
		{
			name: "BaseTime = 2020-03-22 01:02:03, Layout= TimeFormatHour",
			fields: fields{
				timestamp: 1584810123000,
				timezone:  TimezoneTaipei,
				layout:    TimeFormatHour,
			},
			want:    "2020-03-22 01:00:00",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tv := &TimestampValue{
				Timestamp: tt.fields.timestamp,
				Timezone:  tt.fields.timezone,
				Layout:    tt.fields.layout,
			}
			got, err := tv.ToStr()
			if (err != nil) != tt.wantErr {
				t.Errorf("ToStr() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestGetTodayTimestamp(t *testing.T) {
	loc, err := time.LoadLocation(TimezoneTaipei)
	if err != nil {
		t.Log(err)
		return
	}

	type args struct {
		nowTimeBasic time.Time
		timezone     string
		layout       string
	}
	tests := []struct {
		name    string
		args    args
		want    int64
		want1   int64
		wantErr bool
	}{
		{
			name: "today = 2020-01-01 01:01:01",
			args: args{
				nowTimeBasic: time.Date(2020, 1, 1, 1, 1, 1, 1, loc),
				timezone:     TimezoneTaipei,
				layout:       TimeFormatSecond,
			},
			want:    1577808000000,
			want1:   1577894400000,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1, err := GetTodayTimestamp(tt.args.nowTimeBasic, tt.args.timezone, tt.args.layout)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetTodayTimestamp() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			assert.Equal(t, tt.want, got)
			assert.Equal(t, tt.want1, got1)
		})
	}
}

func TestGetYesterdayTimestamp(t *testing.T) {
	type args struct {
		timeStr  string
		timezone string
		layout   string
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
				timeStr:  "2020-03-22 00:00:00",
				timezone: TimezoneTaipei,
				layout:   TimeFormatSecond,
			},
			want:    1584720000000,
			want1:   1584806400000,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1, err := GetYesterdayTimestamp(tt.args.timeStr, tt.args.timezone, tt.args.layout)
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
	loc, err := time.LoadLocation(TimezoneTaipei)
	if err != nil {
		t.Log(err)
		return
	}

	type args struct {
		timeStr  string
		timezone string
		layout   string
		years    int
		months   int
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
				timeStr:  "2020-03-22 00:00:00",
				timezone: TimezoneTaipei,
				layout:   TimeFormatSecond,
				years:    0,
				months:   -2,
			},
			want:    time.Date(2020, 1, 1, 0, 0, 0, 0, loc),
			want1:   time.Date(2020, 1, 31, 0, 0, 0, 0, loc),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1, err := GetMonthStartTimeAndEndTime(tt.args.timeStr, tt.args.timezone, tt.args.layout, tt.args.years, tt.args.months)
			if (err != nil) != tt.wantErr {
				assert.NoError(t, err, "GetMonthStartTimeAndEndTime error = %v", err)
				return
			}

			assert.Equal(t, tt.want, got)
			assert.Equal(t, tt.want1, got1)
		})
	}
}

func TestGetTimeSubDay(t *testing.T) {
	loc, err := time.LoadLocation(TimezoneTaipei)
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
				t1: time.Date(2007, 1, 2, 23, 59, 59, 0, loc),
				t2: time.Date(2007, 1, 3, 0, 0, 0, 0, loc),
			},
			want: 1,
		},
		{
			name: "just one day should return 1, 2007",
			args: args{
				t1: time.Date(2007, 1, 2, 23, 59, 59, 0, loc),
				t2: time.Date(2007, 1, 3, 23, 59, 59, 0, loc),
			},
			want: 1,
		},
		{
			name: "just one day should return 1, 2017",
			args: args{
				t1: time.Date(2017, 9, 1, 10, 0, 0, 0, loc),
				t2: time.Date(2017, 9, 2, 11, 0, 0, 0, loc),
			},
			want: 1,
		},
		{
			name: "just one day should return 2",
			args: args{
				t1: time.Date(2007, 1, 2, 23, 59, 59, 0, loc),
				t2: time.Date(2007, 1, 4, 0, 0, 0, 0, loc),
			},
			want: 2,
		},
		{
			name: "just one day should return 3",
			args: args{
				t1: time.Date(2007, 1, 2, 0, 0, 0, 0, loc),
				t2: time.Date(2007, 1, 5, 0, 0, 0, 0, loc),
			},
			want: 3,
		},
		{
			name: "just one month:31 days should return 31",
			args: args{
				t1: time.Date(2007, 1, 2, 0, 0, 0, 0, loc),
				t2: time.Date(2007, 2, 2, 0, 0, 0, 0, loc),
			},
			want: 31,
		},
		{
			name: "just one month:29 days should return 29",
			args: args{
				t1: time.Date(2000, 2, 1, 0, 0, 0, 0, loc),
				t2: time.Date(2000, 3, 1, 0, 0, 0, 0, loc),
			},
			want: 29,
		},
		{
			name: "just one day: should return 1, loc",
			args: args{
				t1: time.Date(2018, 1, 9, 23, 59, 22, 100, loc),
				t2: time.Date(2018, 1, 10, 0, 0, 1, 100, loc),
			},
			want: 1,
		},
		{
			name: "just one day: should return 1, UTC",
			args: args{
				t1: time.Date(2018, 1, 9, 23, 59, 22, 100, loc),
				t2: time.Date(2018, 1, 10, 0, 0, 1, 100, loc),
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
	loc, err := time.LoadLocation(TimezoneTaipei)
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
				startTime: time.Date(2020, 1, 1, 0, 0, 0, 0, loc),
				t1:        time.Date(2020, 1, 1, 0, 0, 0, 0, loc),
			},
			want: true,
		},
		{
			name: "",
			args: args{
				startTime: time.Date(2020, 1, 1, 0, 0, 0, 0, loc),
				t1:        time.Date(2020, 1, 1, 0, 0, 1, 0, loc),
			},
			want: true,
		},
		{
			name: "",
			args: args{
				startTime: time.Date(2020, 1, 2, 0, 0, 0, 0, loc),
				t1:        time.Date(2020, 1, 1, 23, 59, 59, 0, loc),
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
	loc, err := time.LoadLocation(TimezoneTaipei)
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
				t1: time.Date(2020, 1, 1, 0, 0, 0, 0, loc),
				t2: time.Date(2020, 1, 1, 0, 0, 0, 0, loc),
			},
			want: true,
		},
		{
			name: "",
			args: args{
				t1: time.Date(2020, 1, 1, 0, 0, 0, 0, loc),
				t2: time.Date(2020, 1, 1, 0, 0, 1, 0, loc),
			},
			want: false,
		},
		{
			name: "",
			args: args{
				t1: time.Date(2020, 1, 2, 0, 0, 0, 0, loc),
				t2: time.Date(2020, 1, 1, 23, 59, 59, 0, loc),
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
	loc, err := time.LoadLocation(TimezoneTaipei)
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
				startTime: time.Date(2018, 1, 9, 0, 0, 0, 100, loc),
				endTime:   time.Date(2018, 1, 9, 23, 59, 22, 100, loc),
				t1:        time.Date(2018, 1, 9, 23, 59, 22, 100, loc),
			},
			want: false,
		},
		{
			name: "2018/1/9 00:00:00 ~ 2018/1/9 23:59:22, 2018/1/9 23:59:21",
			args: args{
				startTime: time.Date(2018, 1, 9, 0, 0, 0, 100, loc),
				endTime:   time.Date(2018, 1, 9, 23, 59, 22, 100, loc),
				t1:        time.Date(2018, 1, 9, 23, 59, 21, 100, loc),
			},
			want: true,
		},
		{
			name: "2018/1/9 00:00:00 ~ 2018/1/9 23:59:22, 018/1/9 00:00:00",
			args: args{
				startTime: time.Date(2018, 1, 9, 0, 0, 0, 100, loc),
				endTime:   time.Date(2018, 1, 9, 23, 59, 22, 100, loc),
				t1:        time.Date(2018, 1, 9, 0, 0, 0, 100, loc),
			},
			want: true,
		},
		{
			name: "2018/1/9 00:00:00 ~ 2018/1/9 23:59:22, 2018/1/8 23:59:22",
			args: args{
				startTime: time.Date(2018, 1, 9, 0, 0, 0, 100, loc),
				endTime:   time.Date(2018, 1, 9, 23, 59, 22, 100, loc),
				t1:        time.Date(2018, 1, 8, 23, 59, 22, 100, loc),
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
