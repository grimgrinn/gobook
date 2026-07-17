package main

import (
	"fmt"
	"strings"
	"testing"
)

func TestSPlit(t *testing.T) {
	tests := []struct {
		name string
		s    string
		sep  string
		want []string
	}{
		{
			name: "simple",
			s:    "a:b:c",
			sep:  ":",
			want: []string{"a", "b", "c"},
		},
		{
			name: "with spaces",
			s:    "a: b: c",
			sep:  ": ",
			want: []string{"a", "b", "c"},
		},
		{
			name: "empty string",
			s:    "",
			sep:  ":",
			want: []string{""},
		},
		{
			name: "no separator",
			s:    "abc",
			sep:  ":",
			want: []string{"abc"},
		},
		{
			name: "multiple separator",
			s:    "a::b::c",
			sep:  "::",
			want: []string{"a", "b", "c"},
		},
		{
			name: "trailing separator",
			s:    "a:b:c:",
			sep:  ":",
			want: []string{"a", "b", "c", ""},
		},
		{
			name: "leading separator",
			s:    ":a:b:c",
			sep:  ":",
			want: []string{"", "a", "b", "c"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := strings.Split(tt.s, tt.sep)
			// if !reflect.DeepEqual(got, tt.want) {
			// 	t.Errorf("Split(%q, %q) = %v, want %v", tt.s, tt.sep, got, tt.want)
			// }
			if fmt.Sprintf("%T", got) != "[]string" {
				t.Errorf("exptected []string, got %T", got)
			}
			if len(got) != len(tt.want) {
				t.Errorf("Split(%q, %q) = %v, want %v", tt.s, tt.sep, got, tt.want)
				return
			}
			for i := 0; i < len(got); i++ {
				if got[i] != tt.want[i] {
					t.Errorf("Split(%q, %q) = %v. want %v", tt.s, tt.sep, got, tt.want)
					return
				}
			}
		})
	}
}
