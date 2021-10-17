package tools

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFloatAdd(t *testing.T) {
	type args struct {
		s float64
		p float64
	}
	tests := []struct {
		name    string
		args    args
		want    float64
		wantErr bool
	}{
		{
			name: "test 10.123 + 10.457",
			args: args{
				s: 10.123,
				p: 10.457,
			},
			want:    20.58,
			wantErr: false,
		},
		{
			name: "test 10.1234 + 10.4571",
			args: args{
				s: 10.1234,
				p: 10.4571,
			},
			want:    20.58,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := FloatAdd(tt.args.s, tt.args.p)
			if (err != nil) != tt.wantErr {
				assert.NoError(t, err, "FloatAdd error = %v", err)
				return
			}

			assert.Equal(t, tt.want, got)
		})
	}
}

func TestFloatDiv(t *testing.T) {
	type args struct {
		s float64
		p float64
	}
	tests := []struct {
		name    string
		args    args
		want    float64
		wantErr bool
	}{
		{
			name: "test 1/0",
			args: args{
				s: 1,
				p: 0,
			},
			want:    0,
			wantErr: false,
		},
		{
			name: "test 0/1",
			args: args{
				s: 0,
				p: 1,
			},
			want:    0,
			wantErr: false,
		},
		{
			name: "test 100.123/2",
			args: args{
				s: 100.123,
				p: 2,
			},
			want:    50.0615,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := FloatDiv(tt.args.s, tt.args.p)
			if (err != nil) != tt.wantErr {
				assert.NoError(t, err, "FloatDiv error = %v", err)
				return
			}

			assert.Equal(t, tt.want, got)
		})
	}
}

func TestFloatGT(t *testing.T) {
	type args struct {
		s float64
		p float64
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "test 100.123 > 100.122",
			args: args{
				s: 100.123,
				p: 100.122,
			},
			want: true,
		},
		{
			name: "test 100.122 < 100.123",
			args: args{
				s: 100.122,
				p: 100.123,
			},
			want: false,
		},
		{
			name: "test 100.122 < 100.122",
			args: args{
				s: 100.122,
				p: 100.122,
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := FloatGT(tt.args.s, tt.args.p)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestFloatLT(t *testing.T) {
	type args struct {
		s float64
		p float64
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "test 100.123 > 100.122",
			args: args{
				s: 100.123,
				p: 100.122,
			},
			want: false,
		},
		{
			name: "test 100.122 < 100.123",
			args: args{
				s: 100.122,
				p: 100.123,
			},
			want: true,
		},
		{
			name: "test 100.122 < 100.122",
			args: args{
				s: 100.122,
				p: 100.122,
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := FloatLT(tt.args.s, tt.args.p)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestFloatMul(t *testing.T) {
	type args struct {
		s float64
		p float64
	}
	tests := []struct {
		name    string
		args    args
		want    float64
		wantErr bool
	}{
		{
			name: "test 100.122 * 100.123",
			args: args{
				s: 100.122,
				p: 100.123,
			},
			want:    10024.51,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := FloatMul(tt.args.s, tt.args.p)
			if (err != nil) != tt.wantErr {
				assert.NoError(t, err, "FloatMul error = %v", err)
				return
			}

			assert.Equal(t, tt.want, got)
		})
	}
}

func TestFloatPercent(t *testing.T) {
	type args struct {
		s float64
		p float64
	}
	tests := []struct {
		name    string
		args    args
		want    float64
		wantErr bool
	}{
		{
			name: "test (100.122 / 100.123) * 100",
			args: args{
				s: 100.122,
				p: 100.123,
			},
			want:    99.99,
			wantErr: false,
		},
		{
			name: "test (100.122 / -100.123) * 100",
			args: args{
				s: 100.122,
				p: -100.123,
			},
			want:    0,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := FloatPercent(tt.args.s, tt.args.p)
			if (err != nil) != tt.wantErr {
				assert.NoError(t, err, "FloatPercent error = %v", err)
				return
			}

			assert.Equal(t, tt.want, got)
		})
	}
}

func TestFloatSub(t *testing.T) {
	type args struct {
		s float64
		p float64
	}
	tests := []struct {
		name    string
		args    args
		want    float64
		wantErr bool
	}{
		{
			name: "test 10.123 - 10.457",
			args: args{
				s: 10.123,
				p: 10.457,
			},
			want:    -0.33,
			wantErr: false,
		},
		{
			name: "test 10.1234 - 10.4571",
			args: args{
				s: 10.4571,
				p: 10.1234,
			},
			want:    0.33,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := FloatSub(tt.args.s, tt.args.p)
			if (err != nil) != tt.wantErr {
				assert.NoError(t, err, "FloatSub error = %v", err)
				return
			}

			assert.Equal(t, tt.want, got)
		})
	}
}

func TestFloatTruncate(t *testing.T) {
	type args struct {
		f         float64
		precision int32
	}
	tests := []struct {
		name    string
		args    args
		want    float64
		wantErr bool
	}{
		{
			name: "小數點 第 0 位",
			args: args{
				f:         123.111,
				precision: 0,
			},
			want:    123,
			wantErr: false,
		},
		{
			name: "小數點 第 1 位",
			args: args{
				f:         123.111,
				precision: 1,
			},
			want:    123.1,
			wantErr: false,
		},
		{
			name: "小數點 第 2 位",
			args: args{
				f:         123.111,
				precision: 2,
			},
			want:    123.11,
			wantErr: false,
		},
		{
			name: "小數點 第 3 位",
			args: args{
				f:         123.111,
				precision: 3,
			},
			want:    123.111,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := FloatTruncate(tt.args.f, tt.args.precision)
			if (err != nil) != tt.wantErr {
				assert.NoError(t, err, "FloatTruncate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				assert.Equal(t, tt.want, got)
			}
		})
	}
}

func TestFloatTruncateAndPrecisionSecond(t *testing.T) {
	type args struct {
		f float64
	}
	tests := []struct {
		name    string
		args    args
		want    float64
		wantErr bool
	}{
		{
			name: "小數點 第 2 位",
			args: args{
				f: 123.111,
			},
			want:    123.11,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := FloatTruncateAndPrecisionSecond(tt.args.f)
			if (err != nil) != tt.wantErr {
				assert.NoError(t, err, "FloatTruncateAndPrecisionSecond() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				assert.Equal(t, tt.want, got)
			}
		})
	}
}
