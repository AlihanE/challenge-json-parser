package lexer

import (
	"testing"
)

func TestNew(t *testing.T) {
	type args struct {
		input string
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "1",
			args: args{
				input: "{}",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := New(tt.args.input)
			l.NextToken()
		})
	}
}
