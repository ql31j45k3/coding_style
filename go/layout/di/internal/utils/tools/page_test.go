package tools

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetPageStartAndEnd(t *testing.T) {
	type args struct {
		total int
		page  int
		limit int
	}
	tests := []struct {
		name    string
		args    args
		want    int
		want1   int
		wantErr bool
	}{
		{
			name: "total= 10, page= 12, limit= 2",
			args: args{
				total: 10,
				page:  12,
				limit: 2,
			},
			want:    0,
			want1:   0,
			wantErr: true,
		},
		{
			name: "total= 10, page= 10, limit= 1",
			args: args{
				total: 10,
				page:  10,
				limit: 1,
			},
			want:    9,
			want1:   10,
			wantErr: false,
		},
		{
			name: "total= 10, page= 1, limit= 1",
			args: args{
				total: 10,
				page:  1,
				limit: 1,
			},
			want:    0,
			want1:   1,
			wantErr: false,
		},
		{
			name: "total= 10, page= 2, limit= 2",
			args: args{
				total: 10,
				page:  2,
				limit: 2,
			},
			want:    2,
			want1:   4,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1, err := GetPageStartAndEnd(tt.args.total, tt.args.page, tt.args.limit)
			if (err != nil) != tt.wantErr {
				assert.NoError(t, err, "GetPageStartAndEnd() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			assert.Equal(t, tt.want, got)
			assert.Equal(t, tt.want1, got1)
		})
	}
}
