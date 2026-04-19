// This file is part of PathsHelper library.
// Copyright (C) 2018-2025 Arduino AG (http://www.arduino.cc/)
// Copyright (C) 2021-2026 Alexandre Pujol <alexandre@pujol.io>
// SPDX-License-Identifier: GPL-2.0-only

package paths

import (
	"bytes"
	"path/filepath"
	"strings"
	"testing"
)

var testdataRoot = "../../tests/testdata/paths"

func pathEqualsTo(t *testing.T, expected string, actual *Path) {
	t.Helper()
	got := filepath.ToSlash(actual.String())
	if got != expected {
		t.Errorf("got %v, want %v", got, expected)
	}
}

func makeTestPath(parts []string) *Path {
	p := New(testdataRoot, "fileset")
	for _, part := range parts {
		p = p.Join(part)
	}
	return p
}

func TestFilter(t *testing.T) {
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

func TestPath_IsInsideAnyDir(t *testing.T) {
	tests := []struct {
		name string
		p    string
		dirs []string
		want bool
	}{
		{
			name: "empty dirs",
			p:    "/a/b/c",
			dirs: nil,
			want: false,
		},
		{
			name: "direct child",
			p:    "/a/b/c",
			dirs: []string{"/a/b"},
			want: true,
		},
		{
			name: "nested descendant",
			p:    "/a/b/c/d/e",
			dirs: []string{"/a/b"},
			want: true,
		},
		{
			name: "sibling not under",
			p:    "/a/bc/d",
			dirs: []string{"/a/b"},
			want: false,
		},
		{
			name: "equal path is not under",
			p:    "/a/b",
			dirs: []string{"/a/b"},
			want: false,
		},
		{
			name: "matches one of many",
			p:    "/x/y/z",
			dirs: []string{"/a", "/x", "/p"},
			want: true,
		},
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

func TestPath_New(t *testing.T) {
	tests := []struct {
		name    string
		args    []string
		want    string
		wantNil bool
	}{
		{
			name: "single",
			args: []string{"path"},
			want: "path",
		},
		{
			name: "join",
			args: []string{"path", "path"},
			want: filepath.Join("path", "path"),
		},
		{
			name:    "no-args",
			args:    nil,
			wantNil: true,
		},
		{
			name:    "empty-string",
			args:    []string{""},
			wantNil: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := New(tt.args...)
			if tt.wantNil {
				if got != nil {
					t.Errorf("expected nil, got %v", got)
				}
				return
			}
			if got.String() != tt.want {
				t.Errorf("got %v, want %v", got.String(), tt.want)
			}
		})
	}
}

// testPathCases drives the per-method tests over shared filesystem fixtures.
var testPathCases = []struct {
	name        string
	joinParts   []string // parts joined onto testdataRoot/fileset
	wantPath    string   // forward-slash path string expected
	wantIsDir   bool
	wantExist   bool
	wantExistOK bool // whether ExistCheck should return err==nil
	wantIsDirOK bool // whether IsDirCheck should return err==nil
}{
	{
		name:        "fileset-dir",
		joinParts:   nil,
		wantPath:    testdataRoot + "/fileset",
		wantIsDir:   true,
		wantExist:   true,
		wantExistOK: true,
		wantIsDirOK: true,
	},
	{
		name:        "subdir",
		joinParts:   []string{"folder"},
		wantPath:    testdataRoot + "/fileset/folder",
		wantIsDir:   true,
		wantExist:   true,
		wantExistOK: true,
		wantIsDirOK: true,
	},
	{
		name:        "regular-file",
		joinParts:   []string{"file"},
		wantPath:    testdataRoot + "/fileset/file",
		wantIsDir:   false,
		wantExist:   true,
		wantExistOK: true,
		wantIsDirOK: true,
	},
	{
		name:        "nonexistent",
		joinParts:   []string{"file", "notexistent"},
		wantPath:    testdataRoot + "/fileset/file/notexistent",
		wantIsDir:   false,
		wantExist:   false,
		wantExistOK: true,
		wantIsDirOK: false,
	},
}

func TestPath_Join(t *testing.T) {
	for _, tt := range testPathCases {
		t.Run(tt.name, func(t *testing.T) {
			p := makeTestPath(tt.joinParts)
			pathEqualsTo(t, tt.wantPath, p)
		})
	}
}

func TestPath_IsDirCheck(t *testing.T) {
	for _, tt := range testPathCases {
		t.Run(tt.name, func(t *testing.T) {
			p := makeTestPath(tt.joinParts)
			isDir, err := p.IsDirCheck()
			if isDir != tt.wantIsDir {
				t.Errorf("IsDirCheck() isDir = %v, want %v", isDir, tt.wantIsDir)
			}
			if (err == nil) != tt.wantIsDirOK {
				t.Errorf("IsDirCheck() err = %v, wantOK = %v", err, tt.wantIsDirOK)
			}
		})
	}
}

func TestPath_IsDir(t *testing.T) {
	for _, tt := range testPathCases {
		t.Run(tt.name, func(t *testing.T) {
			p := makeTestPath(tt.joinParts)
			if got := p.IsDir(); got != tt.wantIsDir {
				t.Errorf("IsDir() = %v, want %v", got, tt.wantIsDir)
			}
		})
	}
}

func TestPath_IsNotDir(t *testing.T) {
	for _, tt := range testPathCases {
		t.Run(tt.name, func(t *testing.T) {
			p := makeTestPath(tt.joinParts)
			// IsNotDir returns true only when the path exists and is not a directory
			want := tt.wantExist && !tt.wantIsDir
			if got := p.IsNotDir(); got != want {
				t.Errorf("IsNotDir() = %v, want %v", got, want)
			}
		})
	}
}

func TestPath_ExistCheck(t *testing.T) {
	for _, tt := range testPathCases {
		t.Run(tt.name, func(t *testing.T) {
			p := makeTestPath(tt.joinParts)
			exist, err := p.ExistCheck()
			if exist != tt.wantExist {
				t.Errorf("ExistCheck() exist = %v, want %v", exist, tt.wantExist)
			}
			if (err == nil) != tt.wantExistOK {
				t.Errorf("ExistCheck() err = %v, wantOK = %v", err, tt.wantExistOK)
			}
		})
	}
}

func TestPath_Exist(t *testing.T) {
	for _, tt := range testPathCases {
		t.Run(tt.name, func(t *testing.T) {
			p := makeTestPath(tt.joinParts)
			if got := p.Exist(); got != tt.wantExist {
				t.Errorf("Exist() = %v, want %v", got, tt.wantExist)
			}
		})
	}
}

func TestPath_NotExist(t *testing.T) {
	for _, tt := range testPathCases {
		t.Run(tt.name, func(t *testing.T) {
			p := makeTestPath(tt.joinParts)
			want := !tt.wantExist
			if got := p.NotExist(); got != want {
				t.Errorf("NotExist() = %v, want %v", got, want)
			}
		})
	}
}

func TestPath_ReadDir(t *testing.T) {
	tests := []struct {
		name  string
		parts []string
		want  []string
	}{
		{
			name:  "subfolder",
			parts: []string{"folder"},
			want: []string{
				testdataRoot + "/fileset/folder/.hidden",
				testdataRoot + "/fileset/folder/file2",
				testdataRoot + "/fileset/folder/file3",
				testdataRoot + "/fileset/folder/subfolder",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := makeTestPath(tt.parts)
			list, err := p.ReadDir()
			if err != nil {
				t.Fatal(err)
			}
			if len(list) != len(tt.want) {
				t.Fatalf("got len %d, want %d", len(list), len(tt.want))
			}
			for i, want := range tt.want {
				pathEqualsTo(t, want, list[i])
			}
		})
	}
}

func TestPathList_FilterDirsOnReadDir(t *testing.T) {
	tests := []struct {
		name  string
		parts []string
		want  []string
	}{
		{
			name:  "folder",
			parts: []string{"folder"},
			want: []string{
				testdataRoot + "/fileset/folder/subfolder",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			list, err := makeTestPath(tt.parts).ReadDir()
			if err != nil {
				t.Fatal(err)
			}
			list.FilterDirs()
			if len(list) != len(tt.want) {
				t.Fatalf("got len %d, want %d", len(list), len(tt.want))
			}
			for i, w := range tt.want {
				pathEqualsTo(t, w, list[i])
			}
		})
	}
}

func TestPathList_FilterOutHiddenFilesOnReadDir(t *testing.T) {
	tests := []struct {
		name  string
		parts []string
		want  []string
	}{
		{
			name:  "folder",
			parts: []string{"folder"},
			want: []string{
				testdataRoot + "/fileset/folder/file2",
				testdataRoot + "/fileset/folder/file3",
				testdataRoot + "/fileset/folder/subfolder",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			list, err := makeTestPath(tt.parts).ReadDir()
			if err != nil {
				t.Fatal(err)
			}
			list.FilterOutHiddenFiles()
			if len(list) != len(tt.want) {
				t.Fatalf("got len %d, want %d", len(list), len(tt.want))
			}
			for i, w := range tt.want {
				pathEqualsTo(t, w, list[i])
			}
		})
	}
}

func TestPathList_FilterOutPrefixOnReadDir(t *testing.T) {
	tests := []struct {
		name     string
		parts    []string
		prefixes []string
		want     []string
	}{
		{
			name:     "folder-prefix-file",
			parts:    []string{"folder"},
			prefixes: []string{"file"},
			want: []string{
				testdataRoot + "/fileset/folder/.hidden",
				testdataRoot + "/fileset/folder/subfolder",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			list, err := makeTestPath(tt.parts).ReadDir()
			if err != nil {
				t.Fatal(err)
			}
			list.FilterOutPrefix(tt.prefixes...)
			if len(list) != len(tt.want) {
				t.Fatalf("got len %d, want %d", len(list), len(tt.want))
			}
			for i, w := range tt.want {
				pathEqualsTo(t, w, list[i])
			}
		})
	}
}

func TestPath_FollowSymLink(t *testing.T) {
	tests := []struct {
		name      string
		parts     []string
		entry     string
		wantIsDir bool
	}{
		{
			name:      "symlink-to-folder",
			parts:     nil,
			entry:     "symlinktofolder",
			wantIsDir: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dir := makeTestPath(tt.parts)
			files, err := dir.ReadDir()
			if err != nil {
				t.Fatal(err)
			}
			var match *Path
			for _, file := range files {
				if file.Base() == tt.entry {
					match = file
					break
				}
			}
			if match == nil {
				t.Fatalf("%s not found under %s", tt.entry, dir)
			}
			if err := match.FollowSymLink(); err != nil {
				t.Fatal(err)
			}
			isDir, err := match.IsDirCheck()
			if err != nil {
				t.Fatal(err)
			}
			if isDir != tt.wantIsDir {
				t.Errorf("IsDirCheck() = %v, want %v", isDir, tt.wantIsDir)
			}
		})
	}
}

func TestPath_IsInsideDir(t *testing.T) {
	tests := []struct {
		name string
		a    *Path
		b    *Path
		want bool
	}{
		{
			name: "abs-self",
			a:    New("/a/b/c"),
			b:    New("/a/b/c"),
			want: false,
		},
		{
			name: "abs-parent-inside-child",
			a:    New("/a/b/c"),
			b:    New("/a/b/c/d"),
			want: false,
		},
		{
			name: "abs-child-inside-parent",
			a:    New("/a/b/c/d"),
			b:    New("/a/b/c"),
			want: true,
		},
		{
			name: "abs-deep-parent-inside-child",
			a:    New("/a/b/c"),
			b:    New("/a/b/c/d/e"),
			want: false,
		},
		{
			name: "abs-deep-child-inside-parent",
			a:    New("/a/b/c/d/e"),
			b:    New("/a/b/c"),
			want: true,
		},

		{
			name: "rel-self",
			a:    New("a/b/c"),
			b:    New("a/b/c"),
			want: false,
		},
		{
			name: "rel-parent-inside-child",
			a:    New("a/b/c"),
			b:    New("a/b/c/d"),
			want: false,
		},
		{
			name: "rel-child-inside-parent",
			a:    New("a/b/c/d"),
			b:    New("a/b/c"),
			want: true,
		},
		{
			name: "rel-deep-parent-inside-child",
			a:    New("a/b/c"),
			b:    New("a/b/c/d/e"),
			want: false,
		},
		{
			name: "rel-deep-child-inside-parent",
			a:    New("a/b/c/d/e"),
			b:    New("a/b/c"),
			want: true,
		},
		{
			name: "rel-normalized-inside",
			a:    New("f/../a/b/c/d/e"),
			b:    New("a/b/c"),
			want: true,
		},
		{
			name: "rel-parent-not-inside-normalized",
			a:    New("a/b/c"),
			b:    New("f/../a/b/c/d/e"),
			want: false,
		},
		{
			name: "rel-trailing-dotdot-inside",
			a:    New("a/b/c/d/e/f/.."),
			b:    New("a/b/c"),
			want: true,
		},
		{
			name: "rel-parent-not-inside-trailing-dotdot",
			a:    New("a/b/c"),
			b:    New("a/b/c/d/e/f/.."),
			want: false,
		},

		{
			name: "unrelated-1",
			a:    New("/home/megabug/a15/packages"),
			b:    New("/home/megabug/aide/arduino-1.8.6/hardware/arduino/avr"),
			want: false,
		},
		{
			name: "unrelated-2",
			a:    New("/home/megabug/aide/arduino-1.8.6/hardware/arduino/avr"),
			b:    New("/home/megabug/a15/packages"),
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			isInside, err := tt.a.IsInsideDir(tt.b)
			if err != nil {
				t.Fatal(err)
			}
			if isInside != tt.want {
				t.Errorf("%s IsInsideDir(%s) = %v, want %v", tt.a, tt.b, isInside, tt.want)
			}
		})
	}
}

func TestPath_ReadFileAsLines(t *testing.T) {
	tests := []struct {
		name string
		path []string
		want []string
	}{
		{
			name: "anotherFile",
			path: []string{"fileset", "anotherFile"},
			want: []string{"line 1", "line 2", "", "line 3"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := New(testdataRoot)
			for _, part := range tt.path {
				p = p.Join(part)
			}
			lines, err := p.ReadFileAsLines()
			if err != nil {
				t.Fatal(err)
			}
			if len(lines) != len(tt.want) {
				t.Fatalf("got len %d, want %d", len(lines), len(tt.want))
			}
			for i, want := range tt.want {
				if lines[i] != want {
					t.Errorf("line[%d]: got %v, want %v", i, lines[i], want)
				}
			}
		})
	}
}

func TestPath_CanonicalTempDir(t *testing.T) {
	tests := []struct {
		name string
	}{
		{
			name: "tempdir-canonical",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if TempDir().String() != TempDir().Canonical().String() {
				t.Errorf("got %v, want %v", TempDir().Canonical().String(), TempDir().String())
			}
		})
	}
}

func TestCopyDir(t *testing.T) {
	tests := []struct {
		name string
	}{
		{
			name: "copy-fileset",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tmp, err := MkTempDir("", "")
			if err != nil {
				t.Fatal(err)
			}
			defer tmp.RemoveAll()

			src := New(testdataRoot, "fileset")
			if err := src.CopyDirTo(tmp.Join("dest")); err != nil {
				t.Fatalf("copying dir: %v", err)
			}

			exist, err := tmp.Join("dest", "folder", "subfolder", "file4").ExistCheck()
			if !exist {
				t.Error("expected true")
			}
			if err != nil {
				t.Fatal(err)
			}

			isdir, err := tmp.Join("dest", "folder", "subfolder", "file4").IsDirCheck()
			if isdir {
				t.Error("expected false")
			}
			if err != nil {
				t.Fatal(err)
			}

			if err := src.CopyDirTo(tmp.Join("dest")); err == nil {
				t.Fatal("copying dir to already existing")
			}

			if err := src.Join("file").CopyDirTo(tmp.Join("dest2")); err == nil {
				t.Fatal("copying file as dir")
			}
		})
	}
}

func TestPath_Parents(t *testing.T) {
	tests := []struct {
		name string
		in   string
		want []string
	}{
		{
			name: "absolute",
			in:   "/a/very/long/path",
			want: []string{"/a/very/long/path", "/a/very/long", "/a/very", "/a", "/"},
		},
		{
			name: "relative",
			in:   "a/very/relative/path",
			want: []string{"a/very/relative/path", "a/very/relative", "a/very", "a", "."},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			parents := New(tt.in).Parents()
			if len(parents) != len(tt.want) {
				t.Fatalf("got len %d, want %d", len(parents), len(tt.want))
			}
			for i, want := range tt.want {
				pathEqualsTo(t, want, parents[i])
			}
		})
	}
}

func TestPathList_FilterDirs(t *testing.T) {
	tests := []struct {
		name       string
		parts      []string
		wantBefore []string
		wantAfter  []string
	}{
		{
			name:  "fileset",
			parts: nil,
			wantBefore: []string{
				testdataRoot + "/fileset/anotherFile",
				testdataRoot + "/fileset/file",
				testdataRoot + "/fileset/folder",
				testdataRoot + "/fileset/symlinktofolder",
				testdataRoot + "/fileset/test.txt",
				testdataRoot + "/fileset/test.txt.gz",
			},
			wantAfter: []string{
				testdataRoot + "/fileset/folder",
				testdataRoot + "/fileset/symlinktofolder",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := makeTestPath(tt.parts)
			list, err := p.ReadDir()
			if err != nil {
				t.Fatal(err)
			}
			if len(list) != len(tt.wantBefore) {
				t.Fatalf("got len %d, want %d", len(list), len(tt.wantBefore))
			}
			for i, want := range tt.wantBefore {
				pathEqualsTo(t, want, list[i])
			}
			list.FilterDirs()
			if len(list) != len(tt.wantAfter) {
				t.Fatalf("got len %d, want %d", len(list), len(tt.wantAfter))
			}
			for i, want := range tt.wantAfter {
				pathEqualsTo(t, want, list[i])
			}
		})
	}
}

func TestPathList_FilterOutDirs(t *testing.T) {
	tests := []struct {
		name       string
		readFn     func() (PathList, error)
		wantBefore []string
		wantAfter  []string
	}{
		{
			name: "fileset",
			readFn: func() (PathList, error) {
				return New(testdataRoot, "fileset").ReadDir()
			},
			wantBefore: []string{
				testdataRoot + "/fileset/anotherFile",
				testdataRoot + "/fileset/file",
				testdataRoot + "/fileset/folder",
				testdataRoot + "/fileset/symlinktofolder",
				testdataRoot + "/fileset/test.txt",
				testdataRoot + "/fileset/test.txt.gz",
			},
			wantAfter: []string{
				testdataRoot + "/fileset/anotherFile",
				testdataRoot + "/fileset/file",
				testdataRoot + "/fileset/test.txt",
				testdataRoot + "/fileset/test.txt.gz",
			},
		},
		{
			name: "broken_symlink-dir_1",
			readFn: func() (PathList, error) {
				return New(testdataRoot, "broken_symlink", "dir_1").ReadDirRecursive()
			},
			wantBefore: []string{
				testdataRoot + "/broken_symlink/dir_1/broken_link",
				testdataRoot + "/broken_symlink/dir_1/file2",
				testdataRoot + "/broken_symlink/dir_1/linked_dir",
				testdataRoot + "/broken_symlink/dir_1/linked_dir/file1",
				testdataRoot + "/broken_symlink/dir_1/linked_file",
				testdataRoot + "/broken_symlink/dir_1/real_dir",
				testdataRoot + "/broken_symlink/dir_1/real_dir/file1",
			},
			wantAfter: []string{
				testdataRoot + "/broken_symlink/dir_1/broken_link",
				testdataRoot + "/broken_symlink/dir_1/file2",
				testdataRoot + "/broken_symlink/dir_1/linked_dir/file1",
				testdataRoot + "/broken_symlink/dir_1/linked_file",
				testdataRoot + "/broken_symlink/dir_1/real_dir/file1",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			list, err := tt.readFn()
			if err != nil {
				t.Fatal(err)
			}
			if len(list) != len(tt.wantBefore) {
				t.Fatalf("got len %d, want %d", len(list), len(tt.wantBefore))
			}
			for i, want := range tt.wantBefore {
				pathEqualsTo(t, want, list[i])
			}
			list.FilterOutDirs()
			if len(list) != len(tt.wantAfter) {
				t.Fatalf("got len %d, want %d", len(list), len(tt.wantAfter))
			}
			for i, want := range tt.wantAfter {
				pathEqualsTo(t, want, list[i])
			}
		})
	}
}

func TestPath_EquivalentTo(t *testing.T) {
	wd, err := Getwd()
	if err != nil {
		t.Fatal(err)
	}
	tests := []struct {
		name string
		a    *Path
		b    *Path
		want bool
	}{
		{
			name: "redundant-parent",
			a:    New("file1"),
			b:    New("file1", "somethingelse", ".."),
			want: true,
		},
		{
			name: "redundant-nested",
			a:    New("file1", "abc"),
			b:    New("file1", "abc", "def", ".."),
			want: true,
		},
		{
			name: "abs-vs-relative",
			a:    wd.Join("file1"),
			b:    New("file1"),
			want: true,
		},
		{
			name: "abs-vs-normalized-relative",
			a:    wd.Join("file1"),
			b:    New("file1", "abc", ".."),
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.a.EquivalentTo(tt.b); got != tt.want {
				t.Errorf("EquivalentTo: got %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCanonicalize(t *testing.T) {
	wd, err := Getwd()
	if err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		name string
		in   *Path
		want string
	}{
		{
			name: "existing-file",
			in:   New(testdataRoot, "fileset", "anotherFile"),
			want: wd.Join(testdataRoot, "fileset", "anotherFile").String(),
		},
		{
			name: "nonexistent-file",
			in:   New(testdataRoot, "fileset", "nonexistentFile"),
			want: wd.Join(testdataRoot, "fileset", "nonexistentFile").String(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.in.Canonical()
			if got.String() != tt.want {
				t.Errorf("got %v, want %v", got.String(), tt.want)
			}
		})
	}
}

func TestPath_RelTo(t *testing.T) {
	tests := []struct {
		name    string
		a       *Path
		b       *Path
		want    string
		wantErr bool
	}{
		{
			name: "descendant-to-ancestor",
			a:    New("/my/abs/path/123/456"),
			b:    New("/my/abs/path"),
			want: "../..",
		},
		{
			name: "ancestor-to-descendant",
			a:    New("/my/abs/path"),
			b:    New("/my/abs/path/123/456"),
			want: "123/456",
		},
		{
			name:    "relative-mismatch",
			a:       New("my/path"),
			b:       New("/other/path"),
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res, err := tt.a.RelTo(tt.b)
			if (err != nil) != tt.wantErr {
				t.Fatalf("RelTo() err = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantErr {
				if res != nil {
					t.Errorf("expected nil, got %v", res)
				}
				return
			}
			pathEqualsTo(t, tt.want, res)
		})
	}
}

func TestPath_RelFrom(t *testing.T) {
	tests := []struct {
		name    string
		a       *Path
		b       *Path
		want    string
		wantErr bool
	}{
		{
			name: "descendant-from-ancestor",
			a:    New("/my/abs/path/123/456"),
			b:    New("/my/abs/path"),
			want: "123/456",
		},
		{
			name: "ancestor-from-descendant",
			a:    New("/my/abs/path"),
			b:    New("/my/abs/path/123/456"),
			want: "../..",
		},
		{
			name:    "relative-mismatch",
			a:       New("my/path"),
			b:       New("/other/path"),
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res, err := tt.a.RelFrom(tt.b)
			if (err != nil) != tt.wantErr {
				t.Fatalf("RelFrom() err = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantErr {
				if res != nil {
					t.Errorf("expected nil, got %v", res)
				}
				return
			}
			pathEqualsTo(t, tt.want, res)
		})
	}
}

func TestWriteToTempFile(t *testing.T) {
	tests := []struct {
		name   string
		prefix string
		data   []byte
	}{
		{
			name:   "prefix-test",
			prefix: "prefix",
			data:   []byte("test"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tmpDir := New(testdataRoot, "fileset", "tmp")
			if err := tmpDir.MkdirAll(); err != nil {
				t.Fatal(err)
			}
			defer tmpDir.RemoveAll()

			tmp, err := WriteToTempFile(tt.data, tmpDir, tt.prefix)
			if err != nil {
				t.Fatal(err)
			}
			defer tmp.Remove()

			if !strings.HasPrefix(tmp.Base(), tt.prefix) {
				t.Errorf("base %q does not have prefix %q", tmp.Base(), tt.prefix)
			}
			isInside, err := tmp.IsInsideDir(tmpDir)
			if err != nil {
				t.Fatal(err)
			}
			if !isInside {
				t.Error("expected true")
			}
			data, err := tmp.ReadFile()
			if err != nil {
				t.Fatal(err)
			}
			if !bytes.Equal(data, tt.data) {
				t.Errorf("got %v, want %v", data, tt.data)
			}
		})
	}
}

func TestCopyToSamePath(t *testing.T) {
	tests := []struct {
		name    string
		content []byte
	}{
		{
			name:    "same-file",
			content: []byte("hello"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tmpDir := New(t.TempDir())
			srcFile := tmpDir.Join("test_file")
			dstFile := srcFile

			if err := srcFile.WriteFile(tt.content); err != nil {
				t.Fatal(err)
			}
			content, err := srcFile.ReadFile()
			if err != nil {
				t.Fatal(err)
			}
			if !bytes.Equal(content, tt.content) {
				t.Errorf("got %v, want %v", content, tt.content)
			}

			err = srcFile.CopyTo(dstFile)
			if err == nil {
				t.Fatal("expected error")
			}
			if !strings.Contains(err.Error(), "are the same file") {
				t.Errorf("%q does not contain %q", err.Error(), "are the same file")
			}
		})
	}
}
