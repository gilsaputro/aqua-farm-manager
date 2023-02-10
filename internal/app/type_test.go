package app

import (
	"reflect"
	"testing"
)

func TestUrlID_Int(t *testing.T) {
	tests := []struct {
		name  string
		urlID UrlID
		want  int
	}{
		{
			name:  "get /farm",
			urlID: Farms,
			want:  1,
		},
		{
			name:  "get /pond",
			urlID: Ponds,
			want:  2,
		},
		{
			name:  "get /stat",
			urlID: Stat,
			want:  4,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.urlID.Int(); got != tt.want {
				t.Errorf("UrlID.Int() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUrlID_String(t *testing.T) {
	tests := []struct {
		name  string
		urlID UrlID
		want  string
	}{
		{
			name:  "get /farm",
			urlID: Farms,
			want:  UrlIDName[Farms],
		},
		{
			name:  "get /pond",
			urlID: Ponds,
			want:  UrlIDName[Ponds],
		},
		{
			name:  "get /stat",
			urlID: Stat,
			want:  UrlIDName[Stat],
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.urlID.String(); got != tt.want {
				t.Errorf("UrlID.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUrlID_GetListMethod(t *testing.T) {
	tests := []struct {
		name  string
		urlID UrlID
		want  []string
	}{
		{
			name:  "get /farm",
			urlID: Farms,
			want:  UrlIDMethod[Farms],
		},
		{
			name:  "get /pond",
			urlID: Ponds,
			want:  UrlIDMethod[Ponds],
		},
		{
			name:  "get /stat",
			urlID: Stat,
			want:  UrlIDMethod[Stat],
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.urlID.GetListMethod(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("UrlID.GetListMethod() = %v, want %v", got, tt.want)
			}
		})
	}
}
