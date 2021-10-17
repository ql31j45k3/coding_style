package tools

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSQLArrayToString(t *testing.T) {
	type args struct {
		strs   []string
		column string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "",
			args: args{
				strs:   []string{},
				column: "",
			},
			want: "",
		},
		{
			name: "",
			args: args{
				strs:   []string{"1", "2", "3"},
				column: "column",
			},
			want: " AND (column= ? OR column= ? OR column= ?)",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := SQLArrayToString(tt.args.strs, tt.args.column)
			assert.Equal(t, tt.want, got)
		})
	}
}
