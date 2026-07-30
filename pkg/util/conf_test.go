// apparmor.d - Full set of apparmor profiles
// Copyright (C) 2026 Alexandre Pujol <alexandre@pujol.io>
// SPDX-License-Identifier: GPL-2.0-only

package util

import (
	"os"
	"reflect"
	"slices"
	"testing"

	"github.com/roddhjav/apparmor.d/pkg/paths"
)

// writeConfDir creates a temporary directory seeded with the given files.
func writeConfDir(t *testing.T, files map[string]string) *paths.Path {
	t.Helper()
	dir := paths.New(t.TempDir())
	for name, content := range files {
		if err := dir.Join(name).WriteFile([]byte(content)); err != nil {
			t.Fatalf("write %s: %v", name, err)
		}
	}
	return dir
}

// chmodNoRead removes the read permission from a path, restoring a sane
// mode on cleanup. Skipped for root, which bypasses file modes.
func chmodNoRead(t *testing.T, p *paths.Path) {
	t.Helper()
	if os.Geteuid() == 0 {
		t.Skip("running as root: permission errors cannot be provoked")
	}
	if err := os.Chmod(p.String(), 0o000); err != nil {
		t.Fatalf("chmod %s: %v", p, err)
	}
	t.Cleanup(func() { _ = os.Chmod(p.String(), 0o755) })
}

func TestReadFlagsFile(t *testing.T) {
	tests := []struct {
		name    string
		content string // "" means the file is not created
		want    map[string][]string
	}{
		{
			name:    "missing file",
			content: "",
			want:    map[string][]string{},
		},
		{
			name:    "flags and comments",
			content: "# flags\n\nfoo complain\nbar attach_disconnected,complain\nbaz\n",
			want: map[string][]string{
				"foo": {"complain"},
				"bar": {"attach_disconnected", "complain"},
				"baz": {},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			path := paths.New(t.TempDir()).Join("flags.conf")
			if tt.content != "" {
				if err := path.WriteFile([]byte(tt.content)); err != nil {
					t.Fatalf("write flags: %v", err)
				}
			}
			if got := ReadFlagsFile(path); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ReadFlagsFile() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestReadFlagDirs(t *testing.T) {
	tests := []struct {
		name       string
		vendor     map[string]string // nil means the directory does not exist
		admin      map[string]string
		unreadable bool // remove read permission from the vendor directory
		want       map[string][]string
	}{
		{
			name:       "unreadable dir skipped",
			vendor:     map[string]string{"a.conf": "foo complain\n"},
			admin:      map[string]string{"a.conf": "bar enforce\n"},
			unreadable: true,
			want:       map[string][]string{"bar": {"enforce"}},
		},
		{
			name: "no dirs",
			want: map[string][]string{},
		},
		{
			name:   "vendor only",
			vendor: map[string]string{"a.conf": "foo complain\n"},
			want:   map[string][]string{"foo": {"complain"}},
		},
		{
			name:   "same name admin file replaces vendor file",
			vendor: map[string]string{"a.conf": "foo complain\nbar complain\n"},
			admin:  map[string]string{"a.conf": "foo enforce\n"},
			want:   map[string][]string{"foo": {"enforce"}},
		},
		{
			name:   "different names merge per profile",
			vendor: map[string]string{"00-main.conf": "foo complain\nbar complain\n"},
			admin:  map[string]string{"10-user.conf": "foo enforce\n"},
			want:   map[string][]string{"foo": {"enforce"}, "bar": {"complain"}},
		},
		{
			name:   "later file overrides within a dir",
			vendor: map[string]string{"a.conf": "foo complain\n", "b.conf": "foo enforce\n"},
			want:   map[string][]string{"foo": {"enforce"}},
		},
		{
			name:   "non conf file ignored",
			vendor: map[string]string{"a.conf": "foo complain\n", "b.txt": "bar enforce\n"},
			want:   map[string][]string{"foo": {"complain"}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var dirs []*paths.Path
			for _, files := range []map[string]string{tt.vendor, tt.admin} {
				if files == nil {
					dirs = append(dirs, paths.New(t.TempDir()).Join("missing"))
					continue
				}
				dirs = append(dirs, writeConfDir(t, files))
			}
			if tt.unreadable {
				chmodNoRead(t, dirs[0])
			}
			if got := ReadFlagDirs(dirs...); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ReadFlagDirs() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestReadConfDirs(t *testing.T) {
	tests := []struct {
		name       string
		vendor     map[string]string // nil means the directory does not exist
		admin      map[string]string
		unreadable bool // remove read permission from the vendor directory
		want       []string
	}{
		{
			name:       "unreadable dir skipped",
			vendor:     map[string]string{"a.conf": "foo\n"},
			admin:      map[string]string{"a.conf": "bar\n"},
			unreadable: true,
			want:       []string{"bar"},
		},
		{
			name: "no dirs",
			want: nil,
		},
		{
			name:   "same name admin file replaces vendor file",
			vendor: map[string]string{"a.conf": "# vendor\nfoo\n"},
			admin:  map[string]string{"a.conf": "bar\n"},
			want:   []string{"bar"},
		},
		{
			name:   "different names concatenated",
			vendor: map[string]string{"00-main.conf": "foo\n"},
			admin:  map[string]string{"10-user.conf": "bar\n"},
			want:   []string{"foo", "bar"},
		},
		{
			name:   "files read in alphabetical order",
			vendor: map[string]string{"b.conf": "second\n", "a.conf": "first\n"},
			want:   []string{"first", "second"},
		},
		{
			name:  "missing vendor dir skipped",
			admin: map[string]string{"a.conf": "bar\n"},
			want:  []string{"bar"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var dirs []*paths.Path
			for _, files := range []map[string]string{tt.vendor, tt.admin} {
				if files == nil {
					dirs = append(dirs, paths.New(t.TempDir()).Join("missing"))
					continue
				}
				dirs = append(dirs, writeConfDir(t, files))
			}
			if tt.unreadable {
				chmodNoRead(t, dirs[0])
			}
			if got := ReadConfDirs(dirs...); !slices.Equal(got, tt.want) {
				t.Errorf("ReadConfDirs() = %v, want %v", got, tt.want)
			}
		})
	}
}
