package tool

import (
	"testing"
)

func TestCheckSliceEqual(t *testing.T) {
	type args struct {
		a []string
		b []string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{"test1", args{[]string{"a", "b"}, []string{"b", "a"}}, true},
		{"test2", args{[]string{"a", "b"}, []string{"b", "a", "c"}}, false},
		{"test3", args{[]string{"a", "b", "d"}, []string{"b", "a", "c"}}, false},
		{"test4", args{[]string{"a", "b", "d"}, []string{"b", "a", "c"}}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := CheckSliceEqual(tt.args.a, tt.args.b); got != tt.want {
				t.Errorf("CheckSliceEqual() = %v, want %v", got, tt.want)
			}
		})
	}
}
