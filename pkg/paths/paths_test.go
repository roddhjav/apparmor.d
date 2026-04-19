// apparmor.d - Full set of apparmor profiles
// Copyright (C) 2023-2024 Alexandre Pujol <alexandre@pujol.io>
// SPDX-License-Identifier: GPL-2.0-only

package paths

import "testing"

func Test_Filter(t *testing.T) {
	tests := []struct {
		name string
		src  string
		want string
	}{
		{
			name: "comment",
			src:  "# comment",
			want: "",
		},
		{
			name: "comment with space",
			src:  " # comment",
			want: "",
		},
		{
			name: "no comment",
			src:  "no comment",
			want: "no comment",
		},
		{
			name: "no comment # comment",
			src:  "no comment # comment",
			want: "no comment",
		},
		{
			name: "empty",
			src: `

`,
			want: ``,
		},
		{
			name: "main",
			src: `
# Common profile flags definition for all distributions
# File format: one profile by line using the format: '<profile> <flags>'

bwrap attach_disconnected,mediate_deleted,complain
bwrap-app attach_disconnected,complain

akonadi_akonotes_resource complain # Dev
gnome-disks complain

`,
			want: `bwrap attach_disconnected,mediate_deleted,complain
bwrap-app attach_disconnected,complain
akonadi_akonotes_resource complain
gnome-disks complain
`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotLine := Filter(tt.src)
			if gotLine != tt.want {
				t.Errorf("FilterComment() got = |%v|, want |%v|", gotLine, tt.want)
			}
		})
	}
}

func TestIsInsideAnyDir(t *testing.T) {
	tests := []struct {
		name string
		p    string
		dirs []string
		want bool
	}{
		{name: "empty dirs", p: "/a/b/c", dirs: nil, want: false},
		{name: "direct child", p: "/a/b/c", dirs: []string{"/a/b"}, want: true},
		{name: "nested descendant", p: "/a/b/c/d/e", dirs: []string{"/a/b"}, want: true},
		{name: "sibling not under", p: "/a/bc/d", dirs: []string{"/a/b"}, want: false},
		{name: "equal path is not under", p: "/a/b", dirs: []string{"/a/b"}, want: false},
		{name: "matches one of many", p: "/x/y/z", dirs: []string{"/a", "/x", "/p"}, want: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var dirs []*Path
			for _, d := range tt.dirs {
				dirs = append(dirs, New(d))
			}
			if got := New(tt.p).IsInsideAnyDir(dirs); got != tt.want {
				t.Errorf("IsInsideAnyDir(%q, %v) = %v, want %v", tt.p, tt.dirs, got, tt.want)
			}
		})
	}
}
