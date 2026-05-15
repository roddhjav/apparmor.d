// This file is part of PathsHelper library.
// Copyright (C) 2018-2025 Arduino AG (http://www.arduino.cc/)
// Copyright (C) 2021-2026 Alexandre Pujol <alexandre@pujol.io>
// SPDX-License-Identifier: GPL-2.0-only

package paths

import (
	"fmt"
	"io/fs"
	"os"
	"testing"
	"time"
)

// filesetAll is the expected full list of entries in `fileset` when read recursively and sorted.
var filesetAll = []string{
	testdataRoot + "/fileset/anotherFile",
	testdataRoot + "/fileset/file",
	testdataRoot + "/fileset/folder",
	testdataRoot + "/fileset/folder/.hidden",
	testdataRoot + "/fileset/folder/file2",
	testdataRoot + "/fileset/folder/file3",
	testdataRoot + "/fileset/folder/subfolder",
	testdataRoot + "/fileset/folder/subfolder/file4",
	testdataRoot + "/fileset/symlinktofolder",
	testdataRoot + "/fileset/symlinktofolder/.hidden",
	testdataRoot + "/fileset/symlinktofolder/file2",
	testdataRoot + "/fileset/symlinktofolder/file3",
	testdataRoot + "/fileset/symlinktofolder/subfolder",
	testdataRoot + "/fileset/symlinktofolder/subfolder/file4",
	testdataRoot + "/fileset/test.txt",
	testdataRoot + "/fileset/test.txt.gz",
}

func TestPath_ReadDirRecursive(t *testing.T) {
	tests := []struct {
		name  string
		parts []string
		want  []string
	}{
		{name: "fileset", parts: []string{"fileset"}, want: filesetAll},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := New(testdataRoot)
			for _, part := range tt.parts {
				p = p.Join(part)
			}
			list, err := p.ReadDirRecursive()
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

func TestReadDirRecursiveSymLinkLoop(t *testing.T) {
	tests := []struct {
		name string
		fn   func(tmp *Path) (PathList, error)
	}{
		{
			name: "ReadDirRecursive",
			fn:   func(tmp *Path) (PathList, error) { return tmp.ReadDirRecursive() },
		},
		{
			name: "ReadDirRecursiveFiltered",
			fn:   func(tmp *Path) (PathList, error) { return tmp.ReadDirRecursiveFiltered(nil) },
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tmp, err := MkTempDir("", "")
			if err != nil {
				t.Fatal(err)
			}
			defer tmp.RemoveAll()

			folder := tmp.Join("folder")
			if err := os.Symlink(tmp.String(), folder.String()); err != nil {
				t.Fatal(err)
			}

			l, err := tt.fn(tmp)
			if err == nil {
				t.Fatal("expected error")
			}
			fmt.Println(err)
			if l != nil {
				t.Errorf("expected nil, got %v", l)
			}
		})
	}
}

func TestPath_ReadDirFiltered(t *testing.T) {
	tests := []struct {
		name    string
		filters []ReadDirFilter
		want    []string
	}{
		{
			name:    "no-filter",
			filters: nil,
			want: []string{
				testdataRoot + "/fileset/folder/.hidden",
				testdataRoot + "/fileset/folder/file2",
				testdataRoot + "/fileset/folder/file3",
				testdataRoot + "/fileset/folder/subfolder",
			},
		},
		{
			name:    "only-directories",
			filters: []ReadDirFilter{FilterDirectories()},
			want: []string{
				testdataRoot + "/fileset/folder/subfolder",
			},
		},
		{
			name:    "filter-out-file-prefix",
			filters: []ReadDirFilter{FilterOutPrefixes("file")},
			want: []string{
				testdataRoot + "/fileset/folder/.hidden",
				testdataRoot + "/fileset/folder/subfolder",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			folderPath := New(testdataRoot + "/fileset/folder")
			list, err := folderPath.ReadDir(tt.filters...)
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

func TestPath_ReadDirRecursiveFiltered(t *testing.T) {
	tests := []struct {
		name         string
		recurseFiltr ReadDirFilter
		filters      []ReadDirFilter
		want         []string
	}{
		{
			name:         "no-filters",
			recurseFiltr: nil,
			filters:      nil,
			want:         filesetAll,
		},
		{
			name:         "recurse-only-filter-out-dirs",
			recurseFiltr: FilterOutDirectories(),
			filters:      nil,
			want: []string{
				testdataRoot + "/fileset/anotherFile",
				testdataRoot + "/fileset/file",
				testdataRoot + "/fileset/folder",          // <- listed but not traversed
				testdataRoot + "/fileset/symlinktofolder", // <- listed but not traversed
				testdataRoot + "/fileset/test.txt",
				testdataRoot + "/fileset/test.txt.gz",
			},
		},
		{
			name:         "recurse-nil-filter-out-dirs",
			recurseFiltr: nil,
			filters:      []ReadDirFilter{FilterOutDirectories()},
			want: []string{
				testdataRoot + "/fileset/anotherFile",
				testdataRoot + "/fileset/file",
				testdataRoot + "/fileset/folder/.hidden",
				testdataRoot + "/fileset/folder/file2",
				testdataRoot + "/fileset/folder/file3",
				testdataRoot + "/fileset/folder/subfolder/file4",
				testdataRoot + "/fileset/symlinktofolder/.hidden",
				testdataRoot + "/fileset/symlinktofolder/file2",
				testdataRoot + "/fileset/symlinktofolder/file3",
				testdataRoot + "/fileset/symlinktofolder/subfolder/file4",
				testdataRoot + "/fileset/test.txt",
				testdataRoot + "/fileset/test.txt.gz",
			},
		},
		{
			name:         "both-filter-out-dirs",
			recurseFiltr: FilterOutDirectories(),
			filters:      []ReadDirFilter{FilterOutDirectories()},
			want: []string{
				testdataRoot + "/fileset/anotherFile",
				testdataRoot + "/fileset/file",
				testdataRoot + "/fileset/test.txt",
				testdataRoot + "/fileset/test.txt.gz",
			},
		},
		{
			name:         "recurse-filter-sub-filter-suffix-3",
			recurseFiltr: FilterOutPrefixes("sub"),
			filters:      []ReadDirFilter{FilterOutSuffixes("3")},
			want: []string{
				testdataRoot + "/fileset/anotherFile",
				testdataRoot + "/fileset/file",
				testdataRoot + "/fileset/folder",
				testdataRoot + "/fileset/folder/.hidden",
				testdataRoot + "/fileset/folder/file2",
				testdataRoot + "/fileset/folder/subfolder", // <- subfolder skipped by Prefix("sub")
				testdataRoot + "/fileset/symlinktofolder",
				testdataRoot + "/fileset/symlinktofolder/.hidden",
				testdataRoot + "/fileset/symlinktofolder/file2",
				testdataRoot + "/fileset/symlinktofolder/subfolder", // <- subfolder skipped by Prefix("sub")
				testdataRoot + "/fileset/test.txt",
				testdataRoot + "/fileset/test.txt.gz",
			},
		},
		{
			name:         "recurse-sub-and-filter-suffix3-prefix-fil",
			recurseFiltr: FilterOutPrefixes("sub"),
			filters: []ReadDirFilter{
				AndFilter(FilterOutSuffixes("3"), FilterOutPrefixes("fil")),
			},
			want: []string{
				testdataRoot + "/fileset/anotherFile",
				testdataRoot + "/fileset/folder",
				testdataRoot + "/fileset/folder/.hidden",
				testdataRoot + "/fileset/folder/subfolder",
				testdataRoot + "/fileset/symlinktofolder",
				testdataRoot + "/fileset/symlinktofolder/.hidden",
				testdataRoot + "/fileset/symlinktofolder/subfolder",
				testdataRoot + "/fileset/test.txt",
				testdataRoot + "/fileset/test.txt.gz",
			},
		},
		{
			name:         "recurse-sub-and-filter-suffix3-prefix-fil-suffix-gz",
			recurseFiltr: FilterOutPrefixes("sub"),
			filters: []ReadDirFilter{
				AndFilter(FilterOutSuffixes("3"), FilterOutPrefixes("fil"), FilterOutSuffixes(".gz")),
			},
			want: []string{
				testdataRoot + "/fileset/anotherFile",
				testdataRoot + "/fileset/folder",
				testdataRoot + "/fileset/folder/.hidden",
				testdataRoot + "/fileset/folder/subfolder",
				testdataRoot + "/fileset/symlinktofolder",
				testdataRoot + "/fileset/symlinktofolder/.hidden",
				testdataRoot + "/fileset/symlinktofolder/subfolder",
				testdataRoot + "/fileset/test.txt",
			},
		},
		{
			name:         "or-filter-prefix-sub-or-suffix-tofolder",
			recurseFiltr: OrFilter(FilterPrefixes("sub"), FilterSuffixes("tofolder")),
			filters:      nil,
			want: []string{
				testdataRoot + "/fileset/anotherFile",
				testdataRoot + "/fileset/file",
				testdataRoot + "/fileset/folder",
				testdataRoot + "/fileset/symlinktofolder",
				testdataRoot + "/fileset/symlinktofolder/.hidden",
				testdataRoot + "/fileset/symlinktofolder/file2",
				testdataRoot + "/fileset/symlinktofolder/file3",
				testdataRoot + "/fileset/symlinktofolder/subfolder",
				testdataRoot + "/fileset/symlinktofolder/subfolder/file4",
				testdataRoot + "/fileset/test.txt",
				testdataRoot + "/fileset/test.txt.gz",
			},
		},
		{
			name:         "filter-names-folder",
			recurseFiltr: nil,
			filters:      []ReadDirFilter{FilterNames("folder")},
			want: []string{
				testdataRoot + "/fileset/folder",
			},
		},
		{
			name:         "recurse-symlinktofolder-filter-out-hidden",
			recurseFiltr: FilterNames("symlinktofolder"),
			filters:      []ReadDirFilter{FilterOutNames(".hidden")},
			want: []string{
				testdataRoot + "/fileset/anotherFile",
				testdataRoot + "/fileset/file",
				testdataRoot + "/fileset/folder",
				testdataRoot + "/fileset/symlinktofolder",
				testdataRoot + "/fileset/symlinktofolder/file2",
				testdataRoot + "/fileset/symlinktofolder/file3",
				testdataRoot + "/fileset/symlinktofolder/subfolder",
				testdataRoot + "/fileset/test.txt",
				testdataRoot + "/fileset/test.txt.gz",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testdata := New(testdataRoot, "fileset")
			l, err := testdata.ReadDirRecursiveFiltered(tt.recurseFiltr, tt.filters...)
			if err != nil {
				t.Fatal(err)
			}
			l.Sort()
			if len(l) != len(tt.want) {
				t.Fatalf("got len %d, want %d", len(l), len(tt.want))
			}
			for i, want := range tt.want {
				pathEqualsTo(t, want, l[i])
			}
		})
	}
}

func TestPath_ReadDirRecursiveLoopDetection(t *testing.T) {
	loopsPath := New(testdataRoot, "loops")
	unbuondedReaddir := func(testdir string) (PathList, error) {
		var files PathList
		var err error
		done := make(chan bool)
		go func() {
			files, err = loopsPath.Join(testdir).ReadDirRecursive()
			done <- true
		}()
		deadline := time.After(5 * time.Second)
		tick := time.NewTicker(10 * time.Millisecond)
		defer tick.Stop()
		finished := false
		for !finished {
			select {
			case <-done:
				finished = true
			case <-deadline:
				t.Fatalf("Infinite symlink loop while loading sketch")
			case <-tick.C:
			}
		}
		return files, err
	}

	loopTests := []struct {
		name string
		dir  string
	}{
		{name: "loop_1", dir: "loop_1"},
		{name: "loop_2", dir: "loop_2"},
		{name: "loop_3", dir: "loop_3"},
		{name: "loop_4", dir: "loop_4"},
	}
	for _, tt := range loopTests {
		t.Run(tt.name, func(t *testing.T) {
			l, err := unbuondedReaddir(tt.dir)
			if err == nil || err.Error() != "directories symlink loop detected" {
				t.Fatalf("loop not detected in %s: got %v, want directories symlink loop detected", tt.dir, err)
			}
			if l != nil {
				t.Errorf("expected nil, got %v", l)
			}
		})
	}

	regularTests := []struct {
		name string
		dir  string
		want []string
	}{
		{
			name: "regular_1",
			dir:  "regular_1",
			want: []string{
				testdataRoot + "/loops/regular_1/dir1",
				testdataRoot + "/loops/regular_1/dir1/file1",
				testdataRoot + "/loops/regular_1/dir2",
				testdataRoot + "/loops/regular_1/dir2/file1",
			},
		},
		{
			name: "regular_2",
			dir:  "regular_2",
			want: []string{
				testdataRoot + "/loops/regular_2/dir1",
				testdataRoot + "/loops/regular_2/dir1/file1",
				testdataRoot + "/loops/regular_2/dir2",
				testdataRoot + "/loops/regular_2/dir2/dir1",
				testdataRoot + "/loops/regular_2/dir2/dir1/file1",
				testdataRoot + "/loops/regular_2/dir2/file2",
			},
		},
		{
			name: "regular_3",
			dir:  "regular_3",
			want: []string{
				testdataRoot + "/loops/regular_3/dir1",
				testdataRoot + "/loops/regular_3/dir1/file1",
				testdataRoot + "/loops/regular_3/dir2",
				testdataRoot + "/loops/regular_3/dir2/dir1",
				testdataRoot + "/loops/regular_3/dir2/dir1/file1",
				testdataRoot + "/loops/regular_3/dir2/file2",
				testdataRoot + "/loops/regular_3/link", // broken symlink reported in files
			},
		},
	}
	for _, tt := range regularTests {
		t.Run(tt.name, func(t *testing.T) {
			l, err := unbuondedReaddir(tt.dir)
			if err != nil {
				t.Fatal(err)
			}
			if len(l) != len(tt.want) {
				t.Fatalf("got len %d, want %d", len(l), len(tt.want))
			}
			l.Sort()
			for i, want := range tt.want {
				pathEqualsTo(t, want, l[i])
			}
		})
	}

	t.Run("regular_4_with_permission_error", func(t *testing.T) {
		dir1 := loopsPath.Join("regular_4_with_permission_error", "dir1")

		l, err := unbuondedReaddir("regular_4_with_permission_error")
		if err != nil {
			t.Fatal(err)
		}
		if len(l) == 0 {
			t.Error("expected non-empty list")
		}

		dir1Stat, err := dir1.Stat()
		if err != nil {
			t.Fatal(err)
		}
		if err := dir1.Chmod(fs.FileMode(0)); err != nil {
			t.Fatal(err)
		}
		t.Cleanup(func() {
			dir1.Chmod(dir1Stat.Mode())
		})

		l, err = unbuondedReaddir("regular_4_with_permission_error")
		if err == nil {
			t.Fatal("expected error")
		}
		if l != nil {
			t.Errorf("expected nil, got %v", l)
		}
	})

}
