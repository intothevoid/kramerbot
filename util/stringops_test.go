package util

import "testing"

func TestShortenString(t *testing.T) {
	type args struct {
		str    string
		length int
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{name: "test1", args: args{str: "this is a test string", length: 4}, want: "this"},
		{name: "test2", args: args{str: "this is a test string", length: 14}, want: "this is a test"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ShortenString(tt.args.str, tt.args.length); got != tt.want {
				t.Errorf("ShortenString() = %v, want %v", got, tt.want)
			}
		})
	}
}
