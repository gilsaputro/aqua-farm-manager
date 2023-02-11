package farm

import (
	"aqua-farm-manager/internal/domain/farm"
	"reflect"
	"testing"
)

func TestNewFarmHandler(t *testing.T) {
	type args struct {
		domain  farm.FarmDomain
		options []Option
	}
	tests := []struct {
		name string
		args args
		want *FarmHandler
	}{
		{
			name: "success with setting flow",
			args: args{
				domain:  &farm.Farm{},
				options: []Option{WithTimeoutOptions(10)},
			},
			want: &FarmHandler{
				timeoutInSec: 10,
				domain:       &farm.Farm{},
			},
		},
		{
			name: "success without option flow",
			args: args{
				domain:  &farm.Farm{},
				options: []Option{},
			},
			want: &FarmHandler{
				timeoutInSec: 5,
				domain:       &farm.Farm{},
			},
		},
		{
			name: "success with invalid setting flow",
			args: args{
				domain:  &farm.Farm{},
				options: []Option{WithTimeoutOptions(-1)},
			},
			want: &FarmHandler{
				timeoutInSec: 5,
				domain:       &farm.Farm{},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewFarmHandler(tt.args.domain, tt.args.options...); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewFarmHandler() = %v, want %v", got, tt.want)
			}
		})
	}
}
