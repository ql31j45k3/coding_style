package tools

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsNegativeOne(t *testing.T) {
	type args struct {
		i int
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "test i: 0",
			args: args{
				i: 0,
			},
			want: false,
		},
		{
			name: "test i: -1",
			args: args{
				i: -1,
			},
			want: true,
		},
		{
			name: "test i: 1",
			args: args{
				i: 1,
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := IsNegativeOne(tt.args.i)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestIsNotNegativeOne(t *testing.T) {
	type args struct {
		i int
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "test i: 0",
			args: args{
				i: 0,
			},
			want: true,
		},
		{
			name: "test i: -1",
			args: args{
				i: -1,
			},
			want: false,
		},
		{
			name: "test i: 1",
			args: args{
				i: 1,
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := IsNotNegativeOne(tt.args.i)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestIsNotZero(t *testing.T) {
	type args struct {
		i int
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "test i: 0",
			args: args{
				i: 0,
			},
			want: false,
		},
		{
			name: "test i: -1",
			args: args{
				i: -1,
			},
			want: true,
		},
		{
			name: "test i: 1",
			args: args{
				i: 1,
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := IsNotZero(tt.args.i)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestIsZero(t *testing.T) {
	type args struct {
		i int
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "test i: 0",
			args: args{
				i: 0,
			},
			want: true,
		},
		{
			name: "test i: -1",
			args: args{
				i: -1,
			},
			want: false,
		},
		{
			name: "test i: 1",
			args: args{
				i: 1,
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := IsZero(tt.args.i)
			assert.Equal(t, tt.want, got)
		})
	}
}
