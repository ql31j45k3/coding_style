package tools

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestArrayFind(t *testing.T) {
	type args struct {
		list   []string
		subStr string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "",
			args: args{
				list:   []string{"gf", "gfdm", "ag_fish", "ag_sport"},
				subStr: "ag_fish",
			},
			want: true,
		},
		{
			name: "",
			args: args{
				list:   []string{"gf", "gfdm", "ag_fish", "ag_sport"},
				subStr: "pg",
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ArrayFind(tt.args.list, tt.args.subStr)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestArrayStrUnique(t *testing.T) {
	type args struct {
		list []string
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		{
			name: "",
			args: args{
				list: []string{"gf", "gf", "gfdm", "ag_fish", "ag_sport", "gfdm"},
			},
			want: []string{"gf", "gfdm", "ag_fish", "ag_sport"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ArrayStrUnique(tt.args.list)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestAtoi(t *testing.T) {
	type args struct {
		str          string
		defaultValue int
	}
	tests := []struct {
		name    string
		args    args
		want    int
		wantErr bool
	}{
		{
			name: "test 10 Atoi",
			args: args{
				str:          "10",
				defaultValue: 1,
			},
			want:    10,
			wantErr: false,
		},
		{
			name: "test str Atoi",
			args: args{
				str:          "s",
				defaultValue: 1,
			},
			want:    1,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Atoi(tt.args.str, tt.args.defaultValue)
			if (err != nil) != tt.wantErr {
				assert.NoError(t, err, "Atoi error = %v", err)
				return
			}

			assert.Equal(t, tt.want, got)
		})
	}
}

func TestIsEmpty(t *testing.T) {
	type args struct {
		str string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "test str isEmpty true",
			args: args{
				str: "",
			},
			want: true,
		},
		{
			name: "test str isEmpty false",
			args: args{
				str: "ff",
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := IsEmpty(tt.args.str)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestIsNotEmpty(t *testing.T) {
	type args struct {
		str string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "test str isEmpty false",
			args: args{
				str: "",
			},
			want: false,
		},
		{
			name: "test str isEmpty true",
			args: args{
				str: "ff",
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := IsNotEmpty(tt.args.str)
			assert.Equal(t, tt.want, got)
		})
	}
}
