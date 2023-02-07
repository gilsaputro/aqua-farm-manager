package stat

import (
	"aqua-farm-manager/internal/domain/stat"
	"reflect"
	"testing"
)

func TestNewStatHandler(t *testing.T) {
	type args struct {
		stat    stat.StatDomain
		options []Option
	}
	tests := []struct {
		name string
		args args
		want *StatHandler
	}{
		{
			name: "success without option",
			args: args{
				stat:    &stat.Stat{},
				options: []Option{},
			},
			want: &StatHandler{
				stat:         &stat.Stat{},
				timeoutInSec: 5,
			},
		},
		{
			name: "success wit option",
			args: args{
				stat:    &stat.Stat{},
				options: []Option{WithTimeoutOptions(10)},
			},
			want: &StatHandler{
				stat:         &stat.Stat{},
				timeoutInSec: 10,
			},
		},
		{
			name: "success wit invalid option value",
			args: args{
				stat:    &stat.Stat{},
				options: []Option{WithTimeoutOptions(0)},
			},
			want: &StatHandler{
				stat:         &stat.Stat{},
				timeoutInSec: 5,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewStatHandler(tt.args.stat, tt.args.options...); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewStatHandler() = %v, want %v", got, tt.want)
			}
		})
	}
}
