package model

import "testing"

func TestStatus_Value(t *testing.T) {
	tests := []struct {
		name   string
		status Status
		want   int
	}{
		{
			name:   "Get Active Status",
			status: Active,
			want:   1,
		},
		{
			name:   "Get InActive Status",
			status: Inactive,
			want:   2,
		},
		{
			name:   "Get Uknown Status",
			status: Unknown,
			want:   0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.status.Value(); got != tt.want {
				t.Errorf("Status.Value() = %v, want %v", got, tt.want)
			}
		})
	}
}
