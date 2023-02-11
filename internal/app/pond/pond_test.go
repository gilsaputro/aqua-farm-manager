package pond

import (
	"aqua-farm-manager/internal/domain/pond"
	"reflect"
	"testing"
)

func TestNewPondHandler(t *testing.T) {
	type args struct {
		domain  pond.PondDomain
		options []Option
	}
	tests := []struct {
		name string
		args args
		want *PondHandler
	}{
		{
			name: "success with setting flow",
			args: args{
				domain:  &pond.Pond{},
				options: []Option{WithTimeoutOptions(10)},
			},
			want: &PondHandler{
				timeoutInSec: 10,
				domain:       &pond.Pond{},
			},
		},
		{
			name: "success without option flow",
			args: args{
				domain:  &pond.Pond{},
				options: []Option{},
			},
			want: &PondHandler{
				timeoutInSec: 5,
				domain:       &pond.Pond{},
			},
		},
		{
			name: "success with invalid setting flow",
			args: args{
				domain:  &pond.Pond{},
				options: []Option{WithTimeoutOptions(-1)},
			},
			want: &PondHandler{
				timeoutInSec: 5,
				domain:       &pond.Pond{},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewPondHandler(tt.args.domain, tt.args.options...); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewPondHandler() = %v, want %v", got, tt.want)
			}
		})
	}
}
