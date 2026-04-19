// This file is part of PathsHelper library.
// Copyright (C) 2018-2025 Arduino AG (http://www.arduino.cc/)
// Copyright (C) 2021-2026 Alexandre Pujol <alexandre@pujol.io>
// SPDX-License-Identifier: GPL-2.0-only

package paths

import (
	"fmt"
	"testing"
)

func TestPathList_New(t *testing.T) {
	tests := []struct {
		name string
		args []string
		want string
		len  int
	}{
		{
			name: "empty",
			args: nil,
			want: "[]",
			len:  0,
		},
		{
			name: "single",
			args: []string{"test"},
			want: "[test]",
			len:  1,
		},
		{
			name: "three",
			args: []string{"a", "b", "c"},
			want: "[a b c]",
			len:  3,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			list := NewPathList(tt.args...)
			if len(list) != tt.len {
				t.Fatalf("got len %d, want %d", len(list), tt.len)
			}
			if got := fmt.Sprintf("%s", list); got != tt.want {
				t.Errorf("got %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPathList_Contains(t *testing.T) {
	list := NewPathList("a", "b", "c")
	tests := []struct {
		name string
		path string
		want bool
	}{
		{
			name: "not-present",
			path: "d",
			want: false,
		},
		{
			name: "present",
			path: "a",
			want: true,
		},
		{
			name: "equivalent-not-considered",
			path: "d/../a",
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := list.Contains(New(tt.path)); got != tt.want {
				t.Errorf("got %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPathList_ContainsEquivalentTo(t *testing.T) {
	list := NewPathList("a", "b", "c")
	tests := []struct {
		name string
		path string
		want bool
	}{
		{
			name: "not-present",
			path: "d",
			want: false,
		},
		{
			name: "present",
			path: "a",
			want: true,
		},
		{
			name: "equivalent",
			path: "d/../a",
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := list.ContainsEquivalentTo(New(tt.path)); got != tt.want {
				t.Errorf("got %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPathList_Equals(t *testing.T) {
	base := NewPathList("a", "b", "c")
	tests := []struct {
		name string
		a    PathList
		b    PathList
		want bool
	}{
		{
			name: "clone",
			a:    base,
			b:    base.Clone(),
			want: true,
		},
		{
			name: "different-len",
			a:    base,
			b:    NewPathList("a", "b"),
			want: false,
		},
		{
			name: "different-content",
			a:    base,
			b:    NewPathList("a", "b", "d"),
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.a.Equals(tt.b); got != tt.want {
				t.Errorf("got %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPathList_AddIfMissing(t *testing.T) {
	tests := []struct {
		name  string
		start []string
		add   string
		want  string
	}{
		{
			name:  "add-new",
			start: []string{"a", "b", "c"},
			add:   "d",
			want:  "[a b c d]",
		},
		{
			name:  "skip-existing",
			start: []string{"a", "b", "c", "d"},
			add:   "b",
			want:  "[a b c d]",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			list := NewPathList(tt.start...)
			list.AddIfMissing(New(tt.add))
			if got := fmt.Sprintf("%s", list); got != tt.want {
				t.Errorf("got %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPathList_AddAllMissing(t *testing.T) {
	tests := []struct {
		name  string
		start []string
		add   []string
		want  string
	}{
		{
			name:  "mix-new-and-existing",
			start: []string{"a", "b", "c", "d"},
			add:   []string{"a", "e", "i", "o", "u"},
			want:  "[a b c d e i o u]",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			list := NewPathList(tt.start...)
			list.AddAllMissing(NewPathList(tt.add...))
			if got := fmt.Sprintf("%s", list); got != tt.want {
				t.Errorf("got %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPathList_Sort(t *testing.T) {
	tests := []struct {
		name    string
		entries []string
		before  string
		after   string
	}{
		{
			name: "unsorted",
			entries: []string{
				"pointless", "spare", "carve", "unwieldy", "empty",
				"bow", "tub", "grease", "error", "energetic",
				"depend", "property",
			},
			before: "[pointless spare carve unwieldy empty bow tub grease error energetic depend property]",
			after:  "[bow carve depend empty energetic error grease pointless property spare tub unwieldy]",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			list := NewPathList(tt.entries...)
			if got := fmt.Sprintf("%s", list); got != tt.before {
				t.Errorf("before: got %v, want %v", got, tt.before)
			}
			list.Sort()
			if got := fmt.Sprintf("%s", list); got != tt.after {
				t.Errorf("after: got %v, want %v", got, tt.after)
			}
		})
	}
}

var testFilterList = []string{
	"aaaa",
	"bbbb",
	"cccc",
	"dddd",
	"eeff",
	"aaaa/bbbb",
	"eeee/ffff",
	"gggg/hhhh",
}

func TestPathList_FilterPrefix(t *testing.T) {
	tests := []struct {
		name     string
		prefixes []string
		want     string
	}{
		{
			name:     "single-a",
			prefixes: []string{"a"},
			want:     "[aaaa]",
		},
		{
			name:     "single-b",
			prefixes: []string{"b"},
			want:     "[bbbb aaaa/bbbb]",
		},
		{
			name:     "two-a-b",
			prefixes: []string{"a", "b"},
			want:     "[aaaa bbbb aaaa/bbbb]",
		},
		{
			name:     "no-match",
			prefixes: []string{"test"},
			want:     "[]",
		},
		{
			name:     "empty",
			prefixes: nil,
			want:     "[]",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			list := NewPathList(testFilterList...)
			list.FilterPrefix(tt.prefixes...)
			if got := fmt.Sprintf("%s", list); got != tt.want {
				t.Errorf("got %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPathList_FilterOutPrefix(t *testing.T) {
	tests := []struct {
		name     string
		prefixes []string
		want     string
	}{
		{
			name:     "single-b",
			prefixes: []string{"b"},
			want:     "[aaaa cccc dddd eeff eeee/ffff gggg/hhhh]",
		},
		{
			name:     "multi",
			prefixes: []string{"b", "c", "h"},
			want:     "[aaaa dddd eeff eeee/ffff]",
		},
		{
			name:     "empty",
			prefixes: nil,
			want:     "[aaaa bbbb cccc dddd eeff aaaa/bbbb eeee/ffff gggg/hhhh]",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			list := NewPathList(testFilterList...)
			list.FilterOutPrefix(tt.prefixes...)
			if got := fmt.Sprintf("%s", list); got != tt.want {
				t.Errorf("got %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPathList_FilterSuffix(t *testing.T) {
	tests := []struct {
		name     string
		suffixes []string
		want     string
	}{
		{
			name:     "single-a",
			suffixes: []string{"a"},
			want:     "[aaaa]",
		},
		{
			name:     "two-a-h",
			suffixes: []string{"a", "h"},
			want:     "[aaaa gggg/hhhh]",
		},
		{
			name:     "no-match",
			suffixes: []string{"test"},
			want:     "[]",
		},
		{
			name:     "empty",
			suffixes: nil,
			want:     "[]",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			list := NewPathList(testFilterList...)
			list.FilterSuffix(tt.suffixes...)
			if got := fmt.Sprintf("%s", list); got != tt.want {
				t.Errorf("got %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPathList_FilterOutSuffix(t *testing.T) {
	tests := []struct {
		name     string
		suffixes []string
		want     string
	}{
		{
			name:     "single-a",
			suffixes: []string{"a"},
			want:     "[bbbb cccc dddd eeff aaaa/bbbb eeee/ffff gggg/hhhh]",
		},
		{
			name:     "two-a-h",
			suffixes: []string{"a", "h"},
			want:     "[bbbb cccc dddd eeff aaaa/bbbb eeee/ffff]",
		},
		{
			name:     "no-match",
			suffixes: []string{"test"},
			want:     "[aaaa bbbb cccc dddd eeff aaaa/bbbb eeee/ffff gggg/hhhh]",
		},
		{
			name:     "empty",
			suffixes: nil,
			want:     "[aaaa bbbb cccc dddd eeff aaaa/bbbb eeee/ffff gggg/hhhh]",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			list := NewPathList(testFilterList...)
			list.FilterOutSuffix(tt.suffixes...)
			if got := fmt.Sprintf("%s", list); got != tt.want {
				t.Errorf("got %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPathList_Filter(t *testing.T) {
	tests := []struct {
		name string
		fn   func(p *Path) bool
		want string
	}{
		{
			name: "base-equals-bbbb",
			fn:   func(p *Path) bool { return p.Base() == "bbbb" },
			want: "[bbbb aaaa/bbbb]",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			list := NewPathList(testFilterList...)
			list.Filter(tt.fn)
			if got := fmt.Sprintf("%s", list); got != tt.want {
				t.Errorf("got %v, want %v", got, tt.want)
			}
		})
	}
}
